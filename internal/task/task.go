package task

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/redis"
	"context"
	"errors"
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

func enqueueTask(taskType string, t *asynq.Task, options ...asynq.Option) (*asynq.TaskInfo, error) {
	enqueueOptions := append([]asynq.Option{asynq.Queue(taskType)}, options...)
	info, err := client.Enqueue(t, enqueueOptions...)
	if err == nil {
		prometheus.RecordTaskEnqueued(taskType)
	}
	return info, err
}

func wrapHandler(taskType string, h asynq.HandlerFunc) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		start := time.Now()
		err := h(ctx, t)
		if err == nil {
			recordTaskExecution(ctx, t, model.TaskSuccessStatus, nil, nil)
		} else if shouldRecordFailure(ctx, err) {
			recordTaskExecution(ctx, t, model.TaskFailedStatus, nil, err)
		}
		if err != nil {
			log.Logger.Warningf("task %s failed: %s", taskType, err.Error())
			if errors.Is(err, asynq.SkipRetry) || errors.Is(err, asynq.RevokeTask) {
				log.Logger.Debugf("task %s entered terminal failure without retry", taskType)
			}
		} else {
			if duration := time.Since(start); duration > 5*time.Second {
				log.Logger.Debugf("task %s completed slowly: %s", taskType, duration)
			}
		}
		prometheus.RecordTaskProcessed(taskType, time.Since(start).Seconds(), err == nil)
		return err
	}
}

func Init() {
	client = asynq.NewClientFromRedisClient(redis.RDB)
	inspector = asynq.NewInspectorFromRedisClient(redis.RDB)
	mux = asynq.NewServeMux()
	servers = make([]*asynq.Server, 0)
	addServer := func(queue string, concurrency int) {
		if concurrency <= 0 {
			return
		}
		servers = append(servers, asynq.NewServerFromRedisClient(redis.RDB, newServerConfig(queue, concurrency)))
	}
	addServer(startVictimTaskType, config.Env.AsyncQ.Queues.Victim)
	addServer(stopVictimTaskType, config.Env.AsyncQ.Queues.Victim)
	addServer(loadTrafficTaskType, config.Env.AsyncQ.Queues.Traffic)
	addServer(startGeneratorTaskType, config.Env.AsyncQ.Queues.Generator)
	addServer(stopGeneratorTaskType, config.Env.AsyncQ.Queues.Generator)
	addServer(genAttachmentTaskType, config.Env.AsyncQ.Queues.Attachment)
	addServer(sendEmailTaskType, config.Env.AsyncQ.Queues.Email)
	addServer(webhookTaskType, config.Env.AsyncQ.Queues.Webhook)
	addServer(resizeImageTaskType, config.Env.AsyncQ.Queues.Image)

	mux.HandleFunc(sendEmailTaskType, wrapHandler(sendEmailTaskType, HandleSendEmailTask))
	mux.HandleFunc(startGeneratorTaskType, wrapHandler(startGeneratorTaskType, HandleStartGeneratorTask))
	mux.HandleFunc(stopGeneratorTaskType, wrapHandler(stopGeneratorTaskType, HandleStopGeneratorTask))
	mux.HandleFunc(genAttachmentTaskType, wrapHandler(genAttachmentTaskType, HandleGenAttachmentTask))
	mux.HandleFunc(startVictimTaskType, wrapHandler(startVictimTaskType, HandleStartVictimTask))
	mux.HandleFunc(stopVictimTaskType, wrapHandler(stopVictimTaskType, HandleStopVictimTask))
	mux.HandleFunc(loadTrafficTaskType, wrapHandler(loadTrafficTaskType, HandleLoadTrafficTask))
	mux.HandleFunc(webhookTaskType, wrapHandler(webhookTaskType, HandleWebhookTask))
	mux.HandleFunc(resizeImageTaskType, wrapHandler(resizeImageTaskType, HandleResizeImageTask))
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
		go func() {
			defer wg.Done()
			srv.Shutdown()
		}()
	}
	wg.Wait()
	if inspector != nil {
		if err := inspector.Close(); err != nil {
			log.Logger.Warningf("Failed to close task inspector: %s", err)
		}
		inspector = nil
	}
	if client != nil {
		if err := client.Close(); err != nil {
			log.Logger.Warningf("Failed to close task client: %s", err)
		}
		client = nil
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
