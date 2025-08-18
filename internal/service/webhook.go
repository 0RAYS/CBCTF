package service

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

func CreateWebhook(tx *gorm.DB, form f.CreateWebhookForm) (model.Webhook, bool, string) {
	return db.InitWebhookRepo(tx).Create(db.CreateWebhookOptions{
		Name:       form.Name,
		URL:        form.URL,
		Method:     form.Method,
		Headers:    form.Headers,
		Timeout:    form.Timeout,
		RetryCount: form.RetryCount,
		RetryDelay: form.RetryDelay,
		Events:     form.Events,
	})
}

func UpdateWebhook(tx *gorm.DB, webhook model.Webhook, form f.UpdateWebhookForm) (bool, string) {
	return db.InitWebhookRepo(tx).Update(webhook.ID, db.UpdateWebhookOptions{
		Name:       form.Name,
		URL:        form.URL,
		Method:     form.Method,
		Headers:    form.Headers,
		Timeout:    form.Timeout,
		RetryCount: form.RetryCount,
		RetryDelay: form.RetryDelay,
		On:         form.On,
		Events:     form.Events,
	})
}
