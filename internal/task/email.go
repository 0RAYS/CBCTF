package task

import (
	"CBCTF/internal/email"
	"CBCTF/internal/log"
	"CBCTF/internal/prometheus"
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const SendEmailTaskType = "tasks:email"

type SendEmailPayload struct {
	To    string
	Token string
	ID    string
}

func EnqueueSendEmailTask(to, token, id string) (*asynq.TaskInfo, error) {
	payload, err := json.Marshal(SendEmailPayload{To: to, Token: token, ID: id})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(SendEmailTaskType, payload)
	return client.Enqueue(task, asynq.MaxRetry(3), asynq.Timeout(3*time.Minute))
}

func HandleSendEmailTask(_ context.Context, t *asynq.Task) error {
	var payload SendEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	if err := email.SendVerifyEmail(payload.To, payload.Token, payload.ID); err != nil {
		log.Logger.Warningf("Failed to send mail: %s", err)
		prometheus.IncEmailSentMetrics(false)
	} else {
		prometheus.IncEmailSentMetrics(true)
	}
	return nil
}
