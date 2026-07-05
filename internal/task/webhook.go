package task

import (
	"CBCTF/internal/model"
	"CBCTF/internal/webhook"
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

const webhookTaskType = "tasks:webhook"

type WebhookPayload struct {
	Target model.Webhook
	Event  model.Event
}

func EnqueueWebhookTask(event model.Event, target model.Webhook) (*asynq.TaskInfo, error) {
	payload, err := msgpack.Marshal(WebhookPayload{Event: event, Target: target})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(webhookTaskType, payload)
	return enqueueTask(webhookTaskType, task, asynq.MaxRetry(target.Retry), asynq.Timeout(time.Duration(target.Timeout)*time.Second))
}

func HandleWebhookTask(_ context.Context, task *asynq.Task) error {
	var payload WebhookPayload
	if err := msgpack.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	if err := webhook.SendPayload(payload.Event, payload.Target); err != nil {
		return fmt.Errorf("send webhook failed: event_id=%d webhook_id=%d event_type=%s: %w", payload.Event.ID, payload.Target.ID, payload.Event.Type, err)
	}
	return nil
}
