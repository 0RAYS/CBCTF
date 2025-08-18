package model

import (
	"CBCTF/internal/i18n"
	"time"
)

type Email struct {
	From    string    `json:"from"`
	To      string    `json:"to"`
	Subject string    `json:"subject"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
	Success bool      `json:"success"`
	BasicModel
}

func (e Email) GetModelName() string {
	return "Email"
}

func (e Email) GetVersion() uint {
	return e.Version
}

func (e Email) GetBasicModel() BasicModel {
	return e.BasicModel
}

func (e Email) CreateErrorString() string {
	return i18n.CreateEmailError
}
func (e Email) DeleteErrorString() string {
	return i18n.DeleteEmailError
}

func (e Email) GetErrorString() string {
	return i18n.GetEmailError
}

func (e Email) NotFoundErrorString() string {
	return i18n.EmailNotFound
}

func (e Email) UpdateErrorString() string {
	return i18n.UpdateEmailError
}

func (e Email) GetUniqueKey() []string {
	return []string{"id"}
}
