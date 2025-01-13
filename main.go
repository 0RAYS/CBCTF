package main

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/router"
	"fmt"
)

func init() {
	i18n.Init()
	config.Init()
	log.Init()
	db.Init()
}

func main() {
	ip, port := config.Env.GetString("gin.ip"), config.Env.GetString("gin.port")
	if err := router.Init().Run(fmt.Sprintf("%s:%s", ip, port)); err != nil {
		log.Logger.Panicf("Failed to start: %s", err)
	}
}
