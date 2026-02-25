package task

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/redis"
	"context"
	"strings"
	"time"

	"github.com/hibiken/asynq"
)

var (
	srv    *asynq.Server
	client *asynq.Client
	mux    *asynq.ServeMux
)

func wrapHandler(taskType string, h asynq.HandlerFunc) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		start := time.Now()
		err := h(ctx, t)
		prometheus.RecordTaskProcessed(taskType, time.Since(start).Seconds(), err == nil)
		return err
	}
}

func Init() {
	cfg := asynq.Config{
		Concurrency: config.Env.AsyncQ.Concurrency,
		Logger:      log.Logger.WithField("Type", log.TaskLogType),
	}
	switch strings.ToUpper(config.Env.AsyncQ.Log.Level) {
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

	mux.HandleFunc(SendEmailTaskType, wrapHandler(SendEmailTaskType, HandleSendEmailTask))
	mux.HandleFunc(GenAttachmentTaskType, wrapHandler(GenAttachmentTaskType, HandleGenAttachmentTask))
	mux.HandleFunc(WebhookTaskType, wrapHandler(WebhookTaskType, HandleWebhookTask))
	mux.HandleFunc(ResizeImageTaskType, wrapHandler(ResizeImageTaskType, HandleResizeImageTask))
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
