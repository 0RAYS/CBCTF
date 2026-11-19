package db

import (
	"fmt"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type EmailRepo struct {
	BaseRepo[model.Email]
}

func InitEmailRepo(tx *gorm.DB) *EmailRepo {
	return &EmailRepo{
		DB: tx,
	}
}

func (e *EmailRepo) Create(email model.Email) (model.Email, model.RetVal) {
	if res := e.DB.Model(&model.Email{}).Create(&email); res.Error != nil {
		target := fmt.Sprintf("smtp_id=%d", email.SmtpID)
		log.Logger.Warningf("Failed to create Email: %s (target=%s)", res.Error, target)
		return model.Email{}, model.FailRetVal(i18n.Model.Email.CreateError, "Email", target, res.Error)
	}
	if ret := InitSmtpRepo(e.DB).UpdateStatus(email.SmtpID, email.Success, email.CreatedAt); !ret.OK {
		return model.Email{}, ret
	}
	return email, model.SuccessRetVal()
}
