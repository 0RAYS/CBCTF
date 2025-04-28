package model

import (
	"gorm.io/gorm"
	"time"
)

type Event struct {
	ID        string         `gorm:"primarykey" json:"id"`
	User      User           `json:"-"`
	UserID    uint           `json:"user_id"`
	Team      Team           `json:"-"`
	TeamID    uint           `json:"team_id"`
	Contest   Contest        `json:"-"`
	ContestID uint           `json:"contest_id"`
	Usage     Usage          `json:"-"`
	UsageID   uint           `json:"usage_id"`
	Desc      string         `json:"desc"`
	Type      string         `json:"type"`
	IP        string         `json:"ip"`
	Magic     string         `json:"magic"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Version   uint           `gorm:"default:1" json:"-"`
}
