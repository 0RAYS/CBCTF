package main

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"CBCTF/internel/redis"
	"CBCTF/internel/repo"
	"CBCTF/internel/router"
	"errors"
	"fmt"
	"net/http"
)

func init() {
	config.Init()
	log.Init()
	redis.Init()
	repo.Init()
}

func main() {
	ip, port := config.Env.Gin.Host, config.Env.Gin.Port
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", ip, port),
		Handler: router.Init(),
	}
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Logger.Fatalf("Failed to start: %s", err)
	}
}
