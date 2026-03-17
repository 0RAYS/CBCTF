package task

import (
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/webhook"
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

const webhookTaskType = "tasks:webhook"

type WebhookPayload struct {
	Event  model.Event
	Target model.Webhook
}

func EnqueueWebhookTask(event model.Event, target model.Webhook) (*asynq.TaskInfo, error) {
	payload, err := msgpack.Marshal(WebhookPayload{Event: event, Target: target})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(webhookTaskType, payload)
	info, err := client.Enqueue(task, asynq.Queue(webhookTaskType), asynq.MaxRetry(target.Retry), asynq.Timeout(time.Duration(target.Timeout)*time.Second))
	if err == nil {
		prometheus.RecordTaskEnqueued(webhookTaskType)
	}
	return info, err
}

func HandleWebhookTask(_ context.Context, task *asynq.Task) error {
	var payload WebhookPayload
	if err := msgpack.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return webhook.SendPayload(payload.Event, payload.Target)
}
