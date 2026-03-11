package db

import (
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
)

type SmtpRepo struct {
	BaseRepo[model.Smtp]
}

type CreateSmtpOptions struct {
	Address string
	Host    string
	Port    int
	Pwd     string
	On      bool
}

func (c CreateSmtpOptions) Convert2Model() model.Model {
	return model.Smtp{
		Address: c.Address,
		Host:    c.Host,
		Port:    c.Port,
		Pwd:     c.Pwd,
		On:      c.On,
	}
}

type UpdateSmtpOptions struct {
	Address     *string
	Host        *string
	Port        *int
	Pwd         *string
	On          *bool
	Success     *int64
	SuccessLast *time.Time
	Failure     *int64
	FailureLast *time.Time
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
	if u.On != nil {
		options["on"] = *u.On
	}
	if u.Success != nil {
		options["success"] = *u.Success
	}
	if u.SuccessLast != nil {
		options["success_last"] = *u.SuccessLast
	}
	if u.Failure != nil {
		options["failure"] = *u.Failure
	}
	if u.FailureLast != nil {
		options["failure_last"] = *u.FailureLast
	}
	return options
}

type DiffUpdateSmtpOptions struct {
	Success int64
	Failure int64
}

func (d DiffUpdateSmtpOptions) Convert2Expr() map[string]any {
	options := make(map[string]any)
	if d.Success != 0 {
		options["success"] = gorm.Expr("success + ?", d.Success)
	}
	if d.Failure != 0 {
		options["failure"] = gorm.Expr("failure + ?", d.Failure)
	}
	return options
}

func InitSmtpRepo(tx *gorm.DB) *SmtpRepo {
	return &SmtpRepo{
		BaseRepo: BaseRepo[model.Smtp]{
			DB: tx,
		},
	}
}

func (s *SmtpRepo) UpdateStatus(id uint, success bool, last time.Time) model.RetVal {
	var diffOptions DiffUpdateSmtpOptions
	var options UpdateSmtpOptions
	if success {
		diffOptions = DiffUpdateSmtpOptions{
			Success: 1,
		}
		options = UpdateSmtpOptions{
			SuccessLast: &last,
		}
	} else {
		diffOptions = DiffUpdateSmtpOptions{
			Failure: 1,
		}
		options = UpdateSmtpOptions{
			FailureLast: &last,
		}
	}
	if ret := s.DiffUpdate(id, diffOptions); !ret.OK {
		return ret
	}
	return s.Update(id, options)
}
