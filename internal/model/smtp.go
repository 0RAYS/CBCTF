package model

import (
	"time"
)

type Smtp struct {
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

func (s Smtp) GetUniqueField() []string {
	return []string{"id"}
}

func (s Smtp) GetAllowedQueryFields() []string {
	return []string{"id", "address", "host", "on"}
}
