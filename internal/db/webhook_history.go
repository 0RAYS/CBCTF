package db

import (
	"fmt"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type WebhookHistoryRepo struct {
	BaseRepo[model.WebhookHistory]
}

func InitWebhookHistoryRepo(tx *gorm.DB) *WebhookHistoryRepo {
	return &WebhookHistoryRepo{
		DB: tx,
	}
}

func (w *WebhookHistoryRepo) Create(wh model.WebhookHistory) (model.WebhookHistory, model.RetVal) {
	if res := w.DB.Model(&model.WebhookHistory{}).Create(&wh); res.Error != nil {
		target := fmt.Sprintf("webhook_id=%d", wh.WebhookID)
		log.Logger.Warningf("Failed to create WebhookHistory: %s (target=%s)", res.Error, target)
		return model.WebhookHistory{}, model.FailRetVal(i18n.Model.WebhookHistory.CreateError, model.Name(model.WebhookHistory{}), target, res.Error)
	}
	if ret := InitWebhookRepo(w.DB).UpdateStatus(wh.WebhookID, wh.Success, wh.CreatedAt); !ret.OK {
		return model.WebhookHistory{}, ret
	}
	return wh, model.SuccessRetVal()
}
