package model

import (
	"gorm.io/gorm"
	"time"
)

type Victim struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UsageID   uint           `json:"usage_id"`
	Usage     Usage          `json:"-"`
	TeamID    uint           `json:"team_id"`
	Team      Team           `json:"-"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"-"`
	Pods      []Pod          `json:"-"`
	Start     time.Time      `json:"start"`
	Duration  time.Duration  `json:"duration"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Version   uint           `gorm:"default:1" json:"-"`
}
