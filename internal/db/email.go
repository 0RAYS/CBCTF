package db

import (
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
