package model

import (
	"CBCTF/internal/i18n"
	"time"
)

type Smtp struct {
	Emails      []Email   `json:"-"`
	Address     string    `json:"address"`
	Host        string    `json:"host"`
	Port        int       `json:"port"`
	Pwd         string    `json:"pwd"`
	On          bool      `json:"on"`
	Success     int64     `json:"success"`
	SuccessLast time.Time `json:"success_last"`
	Failure     int64     `json:"failure"`
	FailureLast time.Time `json:"failure_last"`
	BasicModel
}

func (s Smtp) GetModelName() string {
	return "Smtp"
}

func (p Smtp) GetVersion() uint {
	return p.Version
}

func (p Smtp) GetBasicModel() BasicModel {
	return p.BasicModel
}

func (p Smtp) CreateErrorString() string {
	return i18n.CreateSmtpError
}

func (p Smtp) DeleteErrorString() string {
	return i18n.DeleteSmtpError
}

func (p Smtp) GetErrorString() string {
	return i18n.GetSmtpError
}

func (p Smtp) NotFoundErrorString() string {
	return i18n.SmtpNotFound
}

func (p Smtp) UpdateErrorString() string {
	return i18n.UpdateSmtpError
}

func (p Smtp) GetUniqueKey() []string {
	return []string{"id"}
}
