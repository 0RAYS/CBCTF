package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"CBCTF/internal/config"
	"CBCTF/internal/cron"
	"CBCTF/internal/db"
	"CBCTF/internal/email"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/oa"
	"CBCTF/internal/redis"
	"CBCTF/internal/router"
	"CBCTF/internal/sys"
	"CBCTF/internal/task"
	"CBCTF/internal/webhook"
)

var server *http.Server

func run() {
	db.Init()
	redis.Init()
	k8s.Init()
	oa.Init()
	email.Init()
	webhook.Init()
	task.Init()
	cron.Init()

	ip, port := config.Env.Gin.Host, config.Env.Gin.Port
	quit := make(chan os.Signal, 1)
	restart := make(chan os.Signal, 1)
	sys.RegisterStopSignals(quit)
	sys.RegisterRestartSignals(restart)
	redis.InitRateLimiter()
	server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", ip, port),
		Handler:           router.Init(),
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       90 * time.Second,
	}
	// Start synchronously: shutdown must never race a not-yet-started scheduler.
	task.Start()
	cron.Start()
	go func() {
		log.Logger.Infof("Server listening at %s:%d", ip, port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Logger.Fatalf("Failed to start: %s", err)
		}
	}()
	for {
		select {
		case <-restart:
			log.Logger.Info("Restarting server...")
			reboot()
			return
		case <-quit:
			log.Logger.Info("Shutting down server...")
			stop()
			return
		}
	}
}

func reboot() {
	stop()
	run()
}

func stop() {
	shutdownHTTP(server)

	cron.Stop()
	task.Stop()
	k8s.Stop()
	cron.FlushBufferedLogs()
	redis.Stop()
	db.Stop()
}

func shutdownHTTP(s *http.Server) {
	if s == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Logger.Warningf("HTTP drain deadline reached: %s", err)
		// Close cancels remaining requests, including Kubernetes diagnostics/log streams.
		_ = s.Close()
	}
}
