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
	To    string
	Token string
	Kind  EmailKind
}

func enqueueEmailTask(payload SendEmailPayload) (*asynq.TaskInfo, error) {
	data, err := msgpack.Marshal(payload)
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(sendEmailTaskType, data)
	return enqueueTask(sendEmailTaskType, task, asynq.MaxRetry(0), asynq.Timeout(3*time.Minute))
}

func EnqueueSendEmailTask(to, token string) (*asynq.TaskInfo, error) {
	return enqueueEmailTask(SendEmailPayload{Kind: EmailKindVerify, To: to, Token: token})
}

func EnqueueSendResetPasswordEmailTask(to, token string) (*asynq.TaskInfo, error) {
	return enqueueEmailTask(SendEmailPayload{Kind: EmailKindResetPassword, To: to, Token: token})
}

func HandleSendEmailTask(_ context.Context, t *asynq.Task) error {
	var payload SendEmailPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	var err error
	emailKind := "unknown"
	switch payload.Kind {
	case EmailKindVerify:
		emailKind = "verify"
		err = email.SendVerifyEmail(payload.To, payload.Token)
		if err == nil {
			log.Logger.Infof("Verify email sent: to=%s", payload.To)
		}
	case EmailKindResetPassword:
		emailKind = "reset_password"
		err = email.SendResetPasswordEmail(payload.To, payload.Token)
		if err == nil {
			log.Logger.Infof("Reset password email sent: to=%s", payload.To)
		}
	default:
		return fmt.Errorf("unknown email kind: %d", payload.Kind)
	}
	if err != nil {
		prometheus.RecordEmailSent(emailKind, false)
		return err
	}
	prometheus.RecordEmailSent(emailKind, true)
	return nil
}
