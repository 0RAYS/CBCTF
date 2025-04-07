package model

import (
	"gorm.io/gorm"
	"time"
)

type Device struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"-"`
	Magic     string         `json:"magic"`
	Count     int            `json:"count"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Version   uint           `gorm:"default:1" json:"-"`
}
