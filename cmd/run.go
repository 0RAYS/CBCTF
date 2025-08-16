package cmd

import (
	"CBCTF/internal/config"
	"CBCTF/internal/cron"
	"CBCTF/internal/db"
	"CBCTF/internal/email"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/redis"
	"CBCTF/internal/router"
	"CBCTF/internal/task"
	"CBCTF/internal/websocket"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func initialize() {
	config.Init()
	log.Init()
	websocket.Init()
	email.Init()
	redis.Init()
	db.Init()
	k8s.Init()
	task.Init()
	cron.Init()
}

func run() {
	initialize()
	ip, port := config.Env.Gin.Host, config.Env.Gin.Port
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	var server *http.Server
	go func() {
		server = &http.Server{
			Addr:    fmt.Sprintf("%s:%d", ip, port),
			Handler: router.Init(),
		}
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Logger.Fatalf("Failed to start: %s", err)
		}
		log.Logger.Infof("Server started at %s:%d", ip, port)
	}()
	go task.Start()
	go cron.Start()
	<-quit
	log.Logger.Info("Shutting down server...")
	if err := server.Shutdown(context.TODO()); err != nil {
		log.Logger.Fatalf("Failed to shutdown server: %s", err)
	}
	task.Stop()
	cron.Stop()
}
