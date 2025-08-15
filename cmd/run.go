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
	"errors"
	"fmt"
	"net/http"
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

func start() {
	ip, port := config.Env.Gin.Host, config.Env.Gin.Port
	go func() {
		server := &http.Server{
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
}

func run() {
	initialize()
	start()
	select {}
}
