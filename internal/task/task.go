package task

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/redis"
	"context"
	"strings"
	"sync"
	"time"

	"github.com/hibiken/asynq"
)

var (
	client  *asynq.Client
	mux     *asynq.ServeMux
	servers []*asynq.Server
)

const (
	defaultQueueName    = "default"
	victimQueueName     = "victim"
	generatorQueueName  = "generator"
	attachmentQueueName = "attachment"
	emailQueueName      = "email"
	webhookQueueName    = "webhook"
	imageQueueName      = "image"
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
	client = asynq.NewClientFromRedisClient(redis.RDB)
	mux = asynq.NewServeMux()
	servers = make([]*asynq.Server, 0)
	addServer := func(queue string, concurrency int) {
		if concurrency <= 0 {
			return
		}
		servers = append(servers, asynq.NewServerFromRedisClient(redis.RDB, newServerConfig(queue, concurrency)))
	}
	addServer(victimQueueName, config.Env.AsyncQ.Queues.Victim)
	addServer(generatorQueueName, config.Env.AsyncQ.Queues.Generator)
	addServer(attachmentQueueName, config.Env.AsyncQ.Queues.Attachment)
	addServer(emailQueueName, config.Env.AsyncQ.Queues.Email)
	addServer(webhookQueueName, config.Env.AsyncQ.Queues.Webhook)
	addServer(imageQueueName, config.Env.AsyncQ.Queues.Image)
	if len(servers) == 0 {
		servers = append(servers, asynq.NewServerFromRedisClient(redis.RDB, newServerConfig(defaultQueueName, config.Env.AsyncQ.Concurrency)))
	}

	mux.HandleFunc(SendEmailTaskType, wrapHandler(SendEmailTaskType, HandleSendEmailTask))
	mux.HandleFunc(StartGeneratorTaskType, wrapHandler(StartGeneratorTaskType, HandleStartGeneratorTask))
	mux.HandleFunc(StopGeneratorTaskType, wrapHandler(StopGeneratorTaskType, HandleStopGeneratorTask))
	mux.HandleFunc(GenAttachmentTaskType, wrapHandler(GenAttachmentTaskType, HandleGenAttachmentTask))
	mux.HandleFunc(StartVictimTaskType, wrapHandler(StartVictimTaskType, HandleStartVictimTask))
	mux.HandleFunc(StopVictimTaskType, wrapHandler(StopVictimTaskType, HandleStopVictimTask))
	mux.HandleFunc(WebhookTaskType, wrapHandler(WebhookTaskType, HandleWebhookTask))
	mux.HandleFunc(ResizeImageTaskType, wrapHandler(ResizeImageTaskType, HandleResizeImageTask))
}

func Start() {
	for _, srv := range servers {
		if err := srv.Start(mux); err != nil {
			log.Logger.Fatalf("Failed to start task server: %s", err.Error())
		}
	}
	log.Logger.Infof("Task servers started: %d", len(servers))
}

func Stop() {
	var wg sync.WaitGroup
	wg.Add(len(servers))
	for _, srv := range servers {
		go func(srv *asynq.Server) {
			defer wg.Done()
			srv.Shutdown()
		}(srv)
	}
	wg.Wait()
}

func queueForTask(taskType string) string {
	switch taskType {
	case StartVictimTaskType, StopVictimTaskType:
		return victimQueueName
	case StartGeneratorTaskType, StopGeneratorTaskType:
		return generatorQueueName
	case GenAttachmentTaskType:
		return attachmentQueueName
	case SendEmailTaskType:
		return emailQueueName
	case WebhookTaskType:
		return webhookQueueName
	case ResizeImageTaskType:
		return imageQueueName
	default:
		return defaultQueueName
	}
}

func newServerConfig(queue string, concurrency int) asynq.Config {
	cfg := asynq.Config{
		Concurrency: concurrency,
		Logger:      log.Logger.WithField("Type", log.TaskLogType).WithField("Queue", queue),
		Queues: map[string]int{
			queue: 1,
		},
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
	return cfg
}
