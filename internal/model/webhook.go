package model

import (
	"CBCTF/internal/i18n"
	"time"
)

type Webhook struct {
	WebhookHistories []WebhookHistory `json:"webhook_histories"`
	Name             string           `json:"name"`
	URL              string           `json:"url"`
	Method           string           `json:"method"`
	Headers          StringMap        `gorm:"type:json" json:"headers"`
	Timeout          int64            `json:"timeout"`
	RetryCount       int64            `json:"retry_count"`
	RetryDelay       int64            `json:"retry_delay"`
	On               bool             `json:"on"`
	Events           StringList       `gorm:"type:json" json:"events"`
	Success          int64            `json:"success"`
	SuccessLast      time.Time        `json:"success_last"`
	Failure          int64            `json:"failure"`
	FailureLast      time.Time        `json:"failure_last"`
	BasicModel
}

func (w Webhook) GetModelName() string {
	return "Webhook"
}

func (w Webhook) GetVersion() uint {
	return w.Version
}

func (w Webhook) GetBasicModel() BasicModel {
	return w.BasicModel
}

func (w Webhook) CreateErrorString() string {
	return i18n.CreateWebhookError
}

func (w Webhook) DeleteErrorString() string {
	return i18n.DeleteWebhookError
}

func (w Webhook) GetErrorString() string {
	return i18n.GetWebhookError
}

func (w Webhook) NotFoundErrorString() string {
	return i18n.WebhookNotFound
}

func (w Webhook) UpdateErrorString() string {
	return i18n.UpdateWebhookError
}

func (w Webhook) GetUniqueKey() []string {
	return []string{"id"}
}
