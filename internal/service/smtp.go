package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/email"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

func ListSmtps(tx *gorm.DB, form dto.ListModelsForm) ([]model.Smtp, int64, model.RetVal) {
	return db.InitSmtpRepo(tx).List(form.Limit, form.Offset)
}

func CreateSmtp(tx *gorm.DB, form dto.CreateSmtpForm) (model.Smtp, model.RetVal) {
	return db.InitSmtpRepo(tx).Create(db.CreateSmtpOptions{
		Address: form.Address,
		Host:    form.Host,
		Port:    form.Port,
		Pwd:     form.Pwd,
	})
}

func UpdateSmtp(tx *gorm.DB, smtp model.Smtp, form dto.UpdateSmtpForm) (model.Smtp, model.RetVal) {
	if ret := db.InitSmtpRepo(tx).Update(smtp.ID, db.UpdateSmtpOptions{
		Address: form.Address,
		Host:    form.Host,
		Port:    form.Port,
		Pwd:     form.Pwd,
		On:      form.On,
	}); !ret.OK {
		return model.Smtp{}, ret
	}
	newSmtp, ret := db.InitSmtpRepo(tx).GetByID(smtp.ID)
	if !ret.OK {
		return model.Smtp{}, ret
	}
	email.DelSenders(smtp)
	if newSmtp.On {
		email.AddSenders(newSmtp)
	}
	return newSmtp, model.SuccessRetVal()
}

func DeleteSmtp(tx *gorm.DB, smtp model.Smtp) model.RetVal {
	if ret := db.InitSmtpRepo(tx).Delete(smtp.ID); !ret.OK {
		return ret
	}
	email.DelSenders(smtp)
	return model.SuccessRetVal()
}
