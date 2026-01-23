package model

import (
	"CBCTF/internal/i18n"
	"time"
)

type WebhookHistory struct {
	WebhookID uint          `json:"webhook_id"`
	Webhook   Webhook       `json:"-"`
	EventID   uint          `json:"event_id"`
	Event     Event         `json:"-"`
	RespCode  int           `json:"resp_code"`
	Duration  time.Duration `json:"duration"`
	Success   bool          `json:"success"`
	Error     string        `json:"error"`
	BaseModel
}

func (w WebhookHistory) GetModelName() string {
	return "WebhookHistory"
}

func (w WebhookHistory) GetBaseModel() BaseModel {
	return w.BaseModel
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

func (w WebhookHistory) NotFoundString() string {
	return i18n.WebhookHistoryNotFound
}

func (w WebhookHistory) UpdateErrorString() string {
	return i18n.UpdateWebhookHistoryError
}

func (w WebhookHistory) GetUniqueKey() []string {
	return []string{"id"}
}

func (w WebhookHistory) GetAllowedQueryFields() []string {
	return []string{"id", "success"}
}
