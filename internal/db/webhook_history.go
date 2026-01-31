package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
)

type WebhookHistoryRepo struct {
	BaseRepo[model.WebhookHistory]
}

type CreateWebhookHistoryOptions struct {
	WebhookID uint
	EventID   uint
	RespCode  int
	Duration  time.Duration
	Success   bool
	Error     string
}

func (c CreateWebhookHistoryOptions) Convert2Model() model.Model {
	return model.WebhookHistory{
		WebhookID: c.WebhookID,
		EventID:   c.EventID,
		RespCode:  c.RespCode,
		Duration:  c.Duration,
		Success:   c.Success,
		Error:     c.Error,
	}
}

func InitWebhookHistoryRepo(tx *gorm.DB) *WebhookHistoryRepo {
	return &WebhookHistoryRepo{
		BaseRepo: BaseRepo[model.WebhookHistory]{
			DB: tx,
		},
	}
}

func (w *WebhookHistoryRepo) Create(options CreateWebhookHistoryOptions) (model.WebhookHistory, model.RetVal) {
	m := options.Convert2Model().(model.WebhookHistory)
	if res := w.DB.Model(&model.WebhookHistory{}).Create(&m); res.Error != nil {
		log.Logger.Warningf("Failed to create WebhookHistory: %s", res.Error)
		return model.WebhookHistory{}, model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]interface{}{"Model": m.ModelName(), "Error": res.Error.Error()}}
	}
	if ret := InitWebhookRepo(w.DB).UpdateStatus(m.WebhookID, m.Success, m.CreatedAt); !ret.OK {
		return model.WebhookHistory{}, ret
	}
	return m, model.SuccessRetVal()
}
