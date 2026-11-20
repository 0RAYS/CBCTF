package task

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"

	"CBCTF/internal/model"
	"CBCTF/internal/webhook"
)

const webhookTaskType = "tasks:webhook"

type WebhookPayload struct {
	Target model.Webhook
	Event  model.Event
}

func EnqueueWebhookTask(event model.Event, target model.Webhook) (*asynq.TaskInfo, error) {
	return enqueueMsgpackTask(webhookTaskType, WebhookPayload{Event: event, Target: target}, asynq.MaxRetry(target.Retry), asynq.Timeout(time.Duration(target.Timeout)*time.Second))
}

func HandleWebhookTask(_ context.Context, task *asynq.Task) error {
	payload, err := decodeMsgpackTask[WebhookPayload](task)
	if err != nil {
		return err
	}
	if err := webhook.SendPayload(payload.Event, payload.Target); err != nil {
		return fmt.Errorf("send webhook failed: event_id=%d webhook_id=%d event_type=%s: %w", payload.Event.ID, payload.Target.ID, payload.Event.Type, err)
	}
	return nil
}
