package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	Suspect string = "suspect"
	Cheater string = "cheater"
	None    string = "none"

	MagicNotMatch string = "Magic in token and request headers are not matched"
	SameFlag      string = "Flag is same"
)

type Cheat struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"-"`
	TeamID    uint           `json:"team_id"`
	Team      Team           `json:"-"`
	ContestID uint           `json:"contest_id"`
	Contest   Contest        `json:"-"`
	Reason    string         `json:"reason"`
	Type      string         `json:"type"`
	Checked   bool           `json:"checked"`
	Cheated   bool           `json:"cheated"`
	Cheats    Strings        `gorm:"type:json" json:"cheats"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Version   uint           `gorm:"default:1" json:"-"`
}
