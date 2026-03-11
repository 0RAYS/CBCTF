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
	"CBCTF/internal/webhook"
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

// log 与 db 需 Init 两次
// 系统配置主要从数据库中读取, 但数据库连接依赖配置文件
// 初次初始化读取数据库中配置覆盖读取配置文件的值
func preInit() {
	config.Init(configPath)
	log.Init()
	db.Init()
}

func run() {
	log.Init()
	db.Init()
	redis.Init()
	k8s.Init()
	email.Init()
	webhook.Init()
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
	task.Stop()
	cron.Stop()
	redis.Stop()
	db.Stop()
	if err := config.Save(); err != nil {
		log.Logger.Warningf("Failed to save config: %s", err)
	}
	if err := server.Shutdown(context.TODO()); err != nil {
		log.Logger.Warningf("Failed to shutdown server: %s", err)
	}
}

func reboot() {
	time.Sleep(time.Second)
	stop()
	preInit()
	run()
}
