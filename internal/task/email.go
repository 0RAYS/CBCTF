package task

import (
	"CBCTF/internal/email"
	"CBCTF/internal/log"
	"CBCTF/internal/prometheus"
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

const sendEmailTaskType = "tasks:email"

// EmailKind 区分同一队列中不同类型的邮件任务
type EmailKind uint8

const (
	EmailKindVerify        EmailKind = 1
	EmailKindResetPassword EmailKind = 2
)

type SendEmailPayload struct {
	Kind  EmailKind
	To    string
	Token string
	ID    string
}

func enqueueEmailTask(payload SendEmailPayload) (*asynq.TaskInfo, error) {
	data, err := msgpack.Marshal(payload)
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(sendEmailTaskType, data)
	info, err := client.Enqueue(task, asynq.Queue(sendEmailTaskType), asynq.MaxRetry(0), asynq.Timeout(3*time.Minute))
	if err == nil {
		prometheus.RecordTaskEnqueued(sendEmailTaskType)
	}
	return info, err
}

func EnqueueSendEmailTask(to, token, id string) (*asynq.TaskInfo, error) {
	return enqueueEmailTask(SendEmailPayload{Kind: EmailKindVerify, To: to, Token: token, ID: id})
}

func EnqueueSendResetPasswordEmailTask(to, token, id string) (*asynq.TaskInfo, error) {
	return enqueueEmailTask(SendEmailPayload{Kind: EmailKindResetPassword, To: to, Token: token, ID: id})
}

func HandleSendEmailTask(_ context.Context, t *asynq.Task) error {
	var payload SendEmailPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	var err error
	switch payload.Kind {
	case EmailKindVerify:
		err = email.SendVerifyEmail(payload.To, payload.Token, payload.ID)
		if err == nil {
			log.Logger.Infof("Verify email sent: to=%s id=%s", payload.To, payload.ID)
		}
	case EmailKindResetPassword:
		err = email.SendResetPasswordEmail(payload.To, payload.Token, payload.ID)
		if err == nil {
			log.Logger.Infof("Reset password email sent: to=%s id=%s", payload.To, payload.ID)
		}
	default:
		return fmt.Errorf("unknown email kind: %d", payload.Kind)
	}
	if err != nil {
		prometheus.RecordEmailSent(false)
		return err
	}
	prometheus.RecordEmailSent(true)
	return nil
}

