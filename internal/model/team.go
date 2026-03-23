package model

import (
	"time"
)

// Team
// BelongsTo Contest
// ManyToMany User
// HasMany Submission
// HasMany TeamFlag
type Team struct {
	ContestID   uint         `gorm:"index;uniqueIndex:idx_teams_name_contest_active,where:deleted_at IS NULL;not null" json:"contest_id"`
	Contest     Contest      `json:"-"`
	Users       []User       `gorm:"many2many:user_teams;" json:"-"`
	Submissions []Submission `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	TeamFlags   []TeamFlag   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name        string       `gorm:"type:varchar(255);uniqueIndex:idx_teams_name_contest_active,where:deleted_at IS NULL;not null" json:"name"`
	Description string       `json:"description"`
	Captcha     string       `json:"-"`
	Picture     FileURL      `json:"picture"`
	Score       float64      `gorm:"default:0" json:"score"`
	Banned      bool         `gorm:"default:false" json:"banned"`
	Hidden      bool         `gorm:"default:false" json:"hidden"`
	CaptainID   uint         `gorm:"index" json:"captain_id"`
	Rank        int          `gorm:"default:-1" json:"rank"`
	Last        time.Time    `gorm:"default:null" json:"last"`
	BaseModel
}
