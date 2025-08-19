package task

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/redis"

	"github.com/hibiken/asynq"
)

var (
	srv    *asynq.Server
	client *asynq.Client
	mux    *asynq.ServeMux
)

func Init() {
	cfg := asynq.Config{
		Concurrency: config.Env.AsyncQ.Concurrency,
		Logger:      log.Logger.WithField("Type", log.TaskLogType),
	}
	switch config.Env.AsyncQ.Level {
	case "DEBUG":
		cfg.LogLevel = asynq.DebugLevel
	case "INFO":
		cfg.LogLevel = asynq.InfoLevel
	case "WARNING":
		cfg.LogLevel = asynq.WarnLevel
	case "ERROR":
		cfg.LogLevel = asynq.ErrorLevel
	default:
		cfg.LogLevel = asynq.WarnLevel
	}
	srv = asynq.NewServerFromRedisClient(redis.RDB, cfg)
	client = asynq.NewClientFromRedisClient(redis.RDB)
	mux = asynq.NewServeMux()

	mux.HandleFunc(SendEmailTaskType, HandleSendEmailTask)
	mux.HandleFunc(GenAttachmentTaskType, HandleGenAttachmentTask)
	mux.HandleFunc(WebhookTaskType, HandleWebhookTask)
}

func Start() {
	if err := srv.Run(mux); err != nil {
		log.Logger.Fatalf("Failed to start task server: %v", err)
	}
	log.Logger.Info("Task server started")
}

func Stop() {
	srv.Shutdown()
}
