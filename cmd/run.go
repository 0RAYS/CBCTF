package cmd

import (
	"CBCTF/internal/config"
	"CBCTF/internal/cron"
	"CBCTF/internal/db"
	"CBCTF/internal/email"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/redis"
	"CBCTF/internal/router"
	"CBCTF/internal/task"
	"CBCTF/internal/webhook"
	"CBCTF/internal/websocket"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var server *http.Server

func run() {
	db.Init()
	redis.Init()
	k8s.Init()
	email.Init()
	webhook.Init()
	websocket.Init()
	task.Init()
	cron.Init()

	ip, port := config.Env.Gin.Host, config.Env.Gin.Port
	quit := make(chan os.Signal, 1)
	restart := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(restart, syscall.SIGUSR1)
	go func() {
		server = &http.Server{
			Addr:    fmt.Sprintf("%s:%d", ip, port),
			Handler: router.Init(),
		}
		log.Logger.Infof("Server listening at %s:%d", ip, port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Logger.Fatalf("Failed to start: %s", err)
		}
	}()
	go task.Start()
	go cron.Start()
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

func stop() {
	if err := server.Shutdown(context.TODO()); err != nil {
		log.Logger.Fatalf("Failed to shutdown server: %s", err)
	}
	task.Stop()
	cron.Stop()
	redis.Stop()
	db.Stop()
}

func reboot() {
	stop()
	time.Sleep(time.Second)
	i18n.Init()
	config.Init()
	log.Init()
	run()
}
