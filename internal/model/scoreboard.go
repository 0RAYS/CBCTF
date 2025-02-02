package model

import (
	"gorm.io/gorm"
	"time"
)

type Scoreboard struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ContestID uint           `json:"contest_id"`
	TeamID    uint           `json:"team_id"`
	Points    uint           `json:"points"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
