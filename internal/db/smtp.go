package db

import (
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
)

type SmtpRepo struct {
	BasicRepo[model.Smtp]
}

type CreateSmtpOptions struct {
	Address string
	Host    string
	Port    string
	Pwd     string
}

func (c CreateSmtpOptions) Convert2Model() model.Model {
	return model.Smtp{
		Address: c.Address,
		Host:    c.Host,
		Port:    c.Port,
		Pwd:     c.Pwd,
	}
}

type UpdateSmtpOptions struct {
	Address *string
	Host    *string
	Port    *string
	Pwd     *string
	Last    *time.Time
	Count   *int64
}

func (u UpdateSmtpOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Address != nil {
		options["address"] = *u.Address
	}
	if u.Host != nil {
		options["host"] = *u.Host
	}
	if u.Port != nil {
		options["port"] = *u.Port
	}
	if u.Pwd != nil {
		options["pwd"] = *u.Pwd
	}
	if u.Last != nil {
		options["last"] = *u.Last
	}
	if u.Count != nil {
		options["count"] = *u.Count
	}
	return options
}

func InitSmtpRepo(tx *gorm.DB) *SmtpRepo {
	return &SmtpRepo{
		BasicRepo: BasicRepo[model.Smtp]{
			DB: tx,
		},
	}
}
