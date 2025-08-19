package task

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/webhook"
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

const WebhookTaskType = "tasks:webhook"

type WebhookPayload struct {
	Event  model.Event
	Target model.Webhook
}

func EnqueueWebhookTask(event model.Event, target model.Webhook) (*asynq.TaskInfo, error) {
	payload, err := msgpack.Marshal(WebhookPayload{Event: event, Target: target})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(WebhookTaskType, payload)
	return client.Enqueue(task, asynq.MaxRetry(target.RetryCount), asynq.Timeout(time.Duration(target.Timeout)*time.Second))
}

func HandleWebhookTask(_ context.Context, task *asynq.Task) error {
	var payload WebhookPayload
	if err := msgpack.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	if err := webhook.SendPayload(payload.Event, payload.Target); err != nil {
		log.Logger.Warningf("Failed to send webhook payload: %v", err)
		return err
	}
	return nil
}
