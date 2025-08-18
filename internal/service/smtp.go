package service

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

func CreateSmtp(tx *gorm.DB, form f.CreateSmtpForm) (model.Smtp, bool, string) {
	return db.InitSmtpRepo(tx).Create(db.CreateSmtpOptions{
		Address: form.Address,
		Host:    form.Host,
		Port:    form.Port,
		Pwd:     form.Pwd,
	})
}

func UpdateSmtp(tx *gorm.DB, smtp model.Smtp, form f.UpdateSmtpForm) (bool, string) {
	return db.InitSmtpRepo(tx).Update(smtp.ID, db.UpdateSmtpOptions{
		Address: form.Address,
		Host:    form.Host,
		Port:    form.Port,
		Pwd:     form.Pwd,
		On:      form.On,
	})
}
