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
	SuccessLast time.Time `gorm:"default:null" json:"success_last"`
	Failure     int64     `json:"failure"`
	FailureLast time.Time `gorm:"default:null" json:"failure_last"`
	BaseModel
}

func (s Smtp) TableName() string {
	return "smtps"
}

func (s Smtp) ModelName() string {
	return "Smtp"
}

func (s Smtp) GetBaseModel() BaseModel {
	return s.BaseModel
}

func (s Smtp) UniqueFields() []string {
	return []string{"id"}
}

func (s Smtp) QueryFields() []string {
	return []string{"id", "address", "host", "on"}
}
