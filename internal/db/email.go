package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type EmailRepo struct {
	BaseRepo[model.Email]
}

func InitEmailRepo(tx *gorm.DB) *EmailRepo {
	return &EmailRepo{
		BaseRepo: BaseRepo[model.Email]{
			DB: tx,
		},
	}
}

func (e *EmailRepo) Create(email model.Email) (model.Email, model.RetVal) {
	if res := e.DB.Model(&model.Email{}).Create(&email); res.Error != nil {
		log.Logger.Warningf("Failed to create Email: %s", res.Error)
		return model.Email{}, model.RetVal{Msg: i18n.Model.Email.CreateError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	if ret := InitSmtpRepo(e.DB).UpdateStatus(email.SmtpID, email.Success, email.CreatedAt); !ret.OK {
		return model.Email{}, ret
	}
	return email, model.SuccessRetVal()
}
