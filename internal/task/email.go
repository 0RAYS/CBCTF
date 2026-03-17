package task

import (
	"CBCTF/internal/email"
	"CBCTF/internal/prometheus"
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

const sendEmailTaskType = "tasks:email"

type SendEmailPayload struct {
	To    string
	Token string
	ID    string
}

func EnqueueSendEmailTask(to, token, id string) (*asynq.TaskInfo, error) {
	payload, err := msgpack.Marshal(SendEmailPayload{To: to, Token: token, ID: id})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(sendEmailTaskType, payload)
	info, err := client.Enqueue(task, asynq.Queue(sendEmailTaskType), asynq.MaxRetry(0), asynq.Timeout(3*time.Minute))
	if err == nil {
		prometheus.RecordTaskEnqueued(sendEmailTaskType)
	}
	return info, err
}

func HandleSendEmailTask(_ context.Context, t *asynq.Task) error {
	var payload SendEmailPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	if err := email.SendVerifyEmail(payload.To, payload.Token, payload.ID); err != nil {
		prometheus.RecordEmailSent(false)
		return err
	}
	prometheus.RecordEmailSent(true)
	return nil
}
