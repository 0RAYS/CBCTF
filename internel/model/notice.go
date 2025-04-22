package model

import (
	"gorm.io/gorm"
	"time"
)

type Notice struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ContestID uint           `json:"contest_id"`
	Contest   Contest        `json:"-"`
	AdminID   uint           `json:"creator_id"`
	Admin     Admin          `json:"-"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Type      string         `json:"type"`
	CreatedAt time.Time      `json:"created"`
	UpdatedAt time.Time      `json:"updated"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Version   uint           `gorm:"default:1" json:"-"`
}
