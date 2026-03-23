package model

import (
	"time"
)

// WebhookHistory
// BelongsTo Webhook
// BelongsTo Event
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
