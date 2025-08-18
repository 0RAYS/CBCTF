package db

import (
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
)

type WebhookHistoryRepo struct {
	BasicRepo[model.WebhookHistory]
}

type CreateWebhookHistoryOptions struct {
	WebhookID  uint
	EventID    uint
	RespCode   int
	Duration   time.Duration
	Success    bool
	Error      string
	RetryCount int64
	NextRetry  time.Time
}

func (c CreateWebhookHistoryOptions) Convert2Model() model.Model {
	return model.WebhookHistory{
		WebhookID:  c.WebhookID,
		EventID:    c.EventID,
		RespCode:   c.RespCode,
		Duration:   c.Duration,
		Success:    c.Success,
		Error:      c.Error,
		RetryCount: c.RetryCount,
		NextRetry:  c.NextRetry,
	}
}

type UpdateWebhookHistoryOptions struct{}

func (u UpdateWebhookHistoryOptions) Convert2Map() map[string]any {
	return make(map[string]any)
}

func InitWebhookHistoryRepo(tx *gorm.DB) *WebhookHistoryRepo {
	return &WebhookHistoryRepo{
		BasicRepo: BasicRepo[model.WebhookHistory]{
			DB: tx,
		},
	}
}
