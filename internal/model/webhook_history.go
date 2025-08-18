package model

import (
	"CBCTF/internal/i18n"
	"time"
)

type WebhookHistory struct {
	WebhookID  uint          `json:"webhook_id"`
	Webhook    Webhook       `json:"webhook"`
	EventID    uint          `json:"event_id"`
	Event      Event         `json:"event"`
	RespCode   int           `json:"resp_code"`
	Duration   time.Duration `json:"duration"`
	Success    bool          `json:"success"`
	Error      string        `json:"error"`
	RetryCount int64         `json:"retry_count"`
	NextRetry  time.Time     `json:"next_retry"`
	BasicModel
}

func (w WebhookHistory) GetModelName() string {
	return "WebhookHistory"
}

func (w WebhookHistory) GetVersion() uint {
	return w.Version
}

func (w WebhookHistory) GetBasicModel() BasicModel {
	return w.BasicModel
}

func (w WebhookHistory) CreateErrorString() string {
	return i18n.CreateWebhookHistoryError
}

func (w WebhookHistory) DeleteErrorString() string {
	return i18n.DeleteWebhookHistoryError
}

func (w WebhookHistory) GetErrorString() string {
	return i18n.GetWebhookHistoryError
}

func (w WebhookHistory) NotFoundErrorString() string {
	return i18n.WebhookHistoryNotFound
}

func (w WebhookHistory) UpdateErrorString() string {
	return i18n.UpdateWebhookHistoryError
}

func (w WebhookHistory) GetUniqueKey() []string {
	return []string{"id"}
}
