package model

import (
	"gorm.io/gorm"
	"time"
)

type Device struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index:idx_user_id_magic,unique;" json:"user_id"`
	Magic     string         `gorm:"index:idx_user_id_magic,unique;" json:"magic"`
	Count     int            `json:"count"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
