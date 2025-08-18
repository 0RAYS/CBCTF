package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
)

type EmailRepo struct {
	BasicRepo[model.Email]
}

type CreateEmailOptions struct {
	SmtpID  uint
	From    string
	To      string
	Subject string
	Content string
	Time    time.Time
	Success bool
}

func (c CreateEmailOptions) Convert2Model() model.Model {
	return model.Email{
		SmtpID:  c.SmtpID,
		From:    c.From,
		To:      c.To,
		Subject: c.Subject,
		Content: c.Content,
		Time:    c.Time,
		Success: c.Success,
	}
}

type UpdateEmailOptions struct {
}

func (u UpdateEmailOptions) Convert2Map() map[string]any {
	return map[string]any{}
}

func InitEmailRepo(tx *gorm.DB) *EmailRepo {
	return &EmailRepo{
		BasicRepo: BasicRepo[model.Email]{
			DB: tx,
		},
	}
}

func (e *EmailRepo) Create(options CreateEmailOptions) (model.Email, bool, string) {
	m := options.Convert2Model().(model.Email)
	if res := e.DB.Model(&model.Email{}).Create(&m); res.Error != nil {
		log.Logger.Warningf("Failed to create Email: %s", res.Error)
		return model.Email{}, false, i18n.CreateEmailError
	}
	if ok, msg := InitSmtpRepo(e.DB).UpdateStatus(m.ID, m.Success, m.Time); !ok {
		return model.Email{}, false, msg
	}
	return m, true, i18n.Success
}
