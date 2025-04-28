package model

import "time"

type Request struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IP        string    `gorm:"size:45;not null" json:"ip"`
	Time      time.Time `gorm:"not null" json:"time"`
	Method    string    `gorm:"size:10;not null" json:"method"`
	Path      string    `gorm:"size:255;not null" json:"path"`
	URL       string    `gorm:"size:255;not null" json:"url"`
	UserAgent string    `gorm:"size:255;not null" json:"user_agent"`
	Status    int       `gorm:"not null" json:"status"`
	Referer   string    `gorm:"size:255" json:"referer"`
	Magic     string    `json:"magic"`
}
