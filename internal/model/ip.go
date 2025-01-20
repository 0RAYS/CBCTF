package model

import (
	"gorm.io/gorm"
	"time"
)

type IP struct {
	ID        uint           `gorm:"primaryKey"`
	IP        string         `gorm:"size:45;not null" json:"ip"`
	Time      time.Time      `gorm:"not null" json:"time"`
	Method    string         `gorm:"size:10;not null" json:"method"`
	URL       string         `gorm:"size:255;not null" json:"url"`
	UserAgent string         `gorm:"size:255;not null" json:"user_agent"`
	Status    int            `gorm:"not null" json:"status"`
	Referer   string         `gorm:"size:255" json:"referer"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
