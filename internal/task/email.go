package task

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"

	"CBCTF/internal/email"
	"CBCTF/internal/log"
	"CBCTF/internal/prometheus"
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
	return enqueueMsgpackTask(sendEmailTaskType, payload, asynq.MaxRetry(0), asynq.Timeout(3*time.Minute))
}

func EnqueueSendEmailTask(to, token string) (*asynq.TaskInfo, error) {
	return enqueueEmailTask(SendEmailPayload{Kind: EmailKindVerify, To: to, Token: token})
}

func EnqueueSendResetPasswordEmailTask(to, token string) (*asynq.TaskInfo, error) {
	return enqueueEmailTask(SendEmailPayload{Kind: EmailKindResetPassword, To: to, Token: token})
}

func HandleSendEmailTask(_ context.Context, t *asynq.Task) error {
	payload, err := decodeMsgpackTask[SendEmailPayload](t)
	if err != nil {
		return err
	}
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
