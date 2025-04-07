package model

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Password  string         `gorm:"not null" json:"-"`
	Email     string         `gorm:"index:idx_email_deleted,unique;not null" json:"email"`
	Avatar    string         `json:"avatar"`
	Verified  bool           `gorm:"default:false" json:"verified"`
	Notices   []Notice       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Version   uint           `gorm:"default:1" json:"-"`
}
