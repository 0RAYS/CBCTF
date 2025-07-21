package cmd

import (
	"CBCTF/internel/config"
	"CBCTF/internel/cron"
	"CBCTF/internel/email"
	"CBCTF/internel/k8s"
	"CBCTF/internel/log"
	"CBCTF/internel/redis"
	db "CBCTF/internel/repo"
	"CBCTF/internel/router"
	"errors"
	"fmt"
	"net/http"
)

var server *http.Server

func initialize() {
	config.Init()
	log.Init()
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
