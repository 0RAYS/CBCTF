package model

import (
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

func (w WebhookHistory) GetUniqueKey() []string {
	return []string{"id"}
}

func (w WebhookHistory) GetAllowedQueryFields() []string {
	return []string{"id", "success"}
}
