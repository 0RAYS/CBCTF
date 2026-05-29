package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

func ListWebhookHistories(tx *gorm.DB, webhook model.Webhook, form dto.ListModelsForm) ([]model.WebhookHistory, int64, model.RetVal) {
	options := db.GetOptions{
		Preloads: map[string]db.GetOptions{"Webhook": {}, "Event": {}},
		Sort:     []string{"id DESC"},
	}
	if webhook.ID > 0 {
		options.Conditions = map[string]any{"webhook_id": webhook.ID}
	}
	return db.InitWebhookHistoryRepo(tx).List(form.Limit, form.Offset, options)
}
