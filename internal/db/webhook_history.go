package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type WebhookHistoryRepo struct {
	BaseRepo[model.WebhookHistory]
}

func InitWebhookHistoryRepo(tx *gorm.DB) *WebhookHistoryRepo {
	return &WebhookHistoryRepo{
		BaseRepo: BaseRepo[model.WebhookHistory]{
			DB: tx,
		},
	}
}

func (w *WebhookHistoryRepo) Create(wh model.WebhookHistory) (model.WebhookHistory, model.RetVal) {
	if res := w.DB.Model(&model.WebhookHistory{}).Create(&wh); res.Error != nil {
		log.Logger.Warningf("Failed to create WebhookHistory: %s", res.Error)
		return model.WebhookHistory{}, model.RetVal{Msg: i18n.Model.WebhookHistory.CreateError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	if ret := InitWebhookRepo(w.DB).UpdateStatus(wh.WebhookID, wh.Success, wh.CreatedAt); !ret.OK {
		return model.WebhookHistory{}, ret
	}
	return wh, model.SuccessRetVal()
}
