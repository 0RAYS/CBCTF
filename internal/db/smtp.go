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
	Port    int
	Pwd     string
	On      bool
}

func (c CreateSmtpOptions) Convert2Model() model.Model {
	return model.Smtp{
		Address:     c.Address,
		Host:        c.Host,
		Port:        c.Port,
		Pwd:         c.Pwd,
		On:          c.On,
		SuccessLast: time.Now(),
		FailureLast: time.Now(),
	}
}

type UpdateSmtpOptions struct {
	Address     *string
	Host        *string
	Port        *int
	Pwd         *string
	On          *bool
	DiffSuccess int64
	Success     *int64
	SuccessLast *time.Time
	DiffFailure int64
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
	if u.DiffSuccess != 0 {
		options["success"] = gorm.Expr("success + ?", u.DiffSuccess)
	}
	if u.Success != nil {
		options["success"] = *u.Success
	}
	if u.SuccessLast != nil {
		options["success_last"] = *u.SuccessLast
	}
	if u.DiffFailure != 0 {
		options["failure"] = gorm.Expr("failure + ?", u.DiffFailure)
	}
	if u.Failure != nil {
		options["failure"] = *u.Failure
	}
	if u.FailureLast != nil {
		options["failure_last"] = *u.FailureLast
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

func (s *SmtpRepo) UpdateStatus(id uint, success bool, last time.Time) (bool, string) {
	var options UpdateSmtpOptions
	if success {
		options = UpdateSmtpOptions{
			DiffSuccess: 1,
			SuccessLast: &last,
		}
	} else {
		options = UpdateSmtpOptions{
			DiffFailure: 1,
			FailureLast: &last,
		}
	}
	return s.Update(id, options)
}
