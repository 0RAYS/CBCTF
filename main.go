package main

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/redis"
	"CBCTF/internal/router"
	"errors"
	"fmt"
	"net/http"
)

var server *http.Server

func Init() {
	config.Init()
	log.Init()
	redis.Init()
	db.Init()
}

func Start() {
	ip, port := config.Env.Gin.Host, config.Env.Gin.Port
	server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", ip, port),
		Handler: router.Init(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Logger.Panicf("Failed to start: %s", err)
		}
	}()
	log.Logger.Infof("Server started at %s:%d", ip, port)
}

func Stop() {
	if err := server.Close(); err != nil {
		log.Logger.Errorf("Failed to stop: %s", err)
		return
	}
	log.Logger.Info("Server stopped")
	db.Close()
	redis.Close()
}

func Restart() {
	log.Logger.Info("Restarting server")
	Stop()
	Init()
	Start()
}

func main() {
	Init()
	config.Watch(Restart)
	Start()
	select {}
}
