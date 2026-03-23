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
	Success     int64     `gorm:"default:0" json:"success"`
	SuccessLast time.Time `gorm:"default:null" json:"success_last"`
	Failure     int64     `gorm:"default:0" json:"failure"`
	FailureLast time.Time `gorm:"default:null" json:"failure_last"`
	BaseModel
}
