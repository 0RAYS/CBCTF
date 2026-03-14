package model

import (
	"time"
)

// Webhook
// HasMany WebhookHistory
type Webhook struct {
	WebhookHistories []WebhookHistory `json:"-"`
	Name             string           `json:"name"`
	URL              string           `json:"url"`
	Method           string           `json:"method"`
	Headers          StringMap        `gorm:"type:json" json:"headers"`
	Timeout          int64            `json:"timeout"`
	Retry            int              `json:"retry"`
	On               bool             `json:"on"`
	Events           StringList       `gorm:"type:json" json:"events"`
	Success          int64            `gorm:"default:0" json:"success"`
	SuccessLast      time.Time        `gorm:"default:null" json:"success_last"`
	Failure          int64            `gorm:"default:0" json:"failure"`
	FailureLast      time.Time        `gorm:"default:null" json:"failure_last"`
	BaseModel
}

func (w Webhook) TableName() string {
	return "webhooks"
}

func (w Webhook) ModelName() string {
	return "Webhook"
}

func (w Webhook) GetBaseModel() BaseModel {
	return w.BaseModel
}

func (w Webhook) UniqueFields() []string {
	return []string{"id"}
}

func (w Webhook) QueryFields() []string {
	return []string{"id", "name", "url", "on", "method"}
}
