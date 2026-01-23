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
	BaseModel
}

func (s Smtp) GetModelName() string {
	return "Smtp"
}

func (s Smtp) GetBaseModel() BaseModel {
	return s.BaseModel
}

func (s Smtp) CreateErrorString() string {
	return i18n.CreateSmtpError
}

func (s Smtp) DeleteErrorString() string {
	return i18n.DeleteSmtpError
}

func (s Smtp) GetErrorString() string {
	return i18n.GetSmtpError
}

func (s Smtp) NotFoundErrorString() string {
	return i18n.SmtpNotFound
}

func (s Smtp) UpdateErrorString() string {
	return i18n.UpdateSmtpError
}

func (s Smtp) GetUniqueKey() []string {
	return []string{"id"}
}

func (s Smtp) GetAllowedQueryFields() []string {
	return []string{"id", "address", "host", "on"}
}
