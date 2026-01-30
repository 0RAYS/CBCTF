package model

import (
	"time"
)

// Team
// BelongsTo Contest
// ManyToMany User
// ManyToMany Submission
// HasMany TeamFlag
type Team struct {
	ContestID   uint         `gorm:"index:idx_name_contest,unique;not null" json:"contest_id"`
	Users       []User       `gorm:"many2many:user_teams;" json:"-"`
	Submissions []Submission `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	TeamFlags   []TeamFlag   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name        string       `gorm:"type:varchar(255);index:idx_name_contest,unique;not null" json:"name"`
	Description string       `json:"description"`
	Captcha     string       `json:"-"`
	Picture     FileURL      `json:"picture"`
	Score       float64      `gorm:"default:0" json:"score"`
	Banned      bool         `gorm:"default:false" json:"banned"`
	Hidden      bool         `gorm:"default:false" json:"hidden"`
	CaptainID   uint         `json:"captain_id"`
	Rank        int          `gorm:"default:-1" json:"rank"`
	Last        time.Time    `gorm:"default:null" json:"last"`
	BaseModel
}

func (t Team) GetModelName() string {
	return "Team"
}

func (t Team) GetBaseModel() BaseModel {
	return t.BaseModel
}

func (t Team) GetUniqueKey() []string {
	return []string{"id"}
}

func (t Team) GetAllowedQueryFields() []string {
	return []string{"id", "name", "description", "banned", "hidden"}
}
