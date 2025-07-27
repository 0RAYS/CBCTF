package cmd

import (
	"CBCTF/internal/config"
	"CBCTF/internal/cron"
	"CBCTF/internal/email"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/redis"
	db "CBCTF/internal/repo"
	"CBCTF/internal/router"
	"CBCTF/internal/websocket"
	"errors"
	"fmt"
	"net/http"
)

var server *http.Server

func initialize() {
	config.Init()
	log.Init()
	websocket.Init()
	email.Init()
	redis.Init()
	db.Init()
	k8s.Init()
	go cron.Init()
}

func start() {

	ip, port := config.Env.Gin.Host, config.Env.Gin.Port
	server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", ip, port),
		Handler: router.Init(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Logger.Fatalf("Failed to start: %s", err)
		}
	}()
	log.Logger.Infof("Server started at %s:%d", ip, port)
	cron.Start()
}

func run() {
	initialize()
	start()
	select {}
}
