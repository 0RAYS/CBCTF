package model

import (
	"gorm.io/gorm"
	"time"
)

type Answer struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TeamID    uint           `json:"team_id"`
	Team      Team           `json:"-"`
	FlagID    uint           `json:"flag_id"`
	Flag      Flag           `json:"-"`
	Value     string         `gorm:"not null" json:"value"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Version   uint           `gorm:"default:1" json:"-"`
}
