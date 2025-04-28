package model

import (
	"gorm.io/gorm"
	"time"
)

type Usage struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ContestID   uint           `json:"contest_id"`
	Contest     Contest        `json:"-"`
	ChallengeID string         `json:"challenge_id"`
	Challenge   Challenge      `json:"-"`
	Flags       []Flag         `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Victims     []Victim       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions []Submission   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Events      []Event        `json:"-"`
	Name        string         `gorm:"not null" json:"name"`
	Desc        string         `json:"desc"`
	Hidden      bool           `json:"hidden"`
	Attempt     int64          `json:"attempt"`
	Dockers     Dockers        `gorm:"type:json" json:"dockers"`
	Hints       Strings        `gorm:"type:json" json:"hints"`
	Tags        Strings        `gorm:"type:json" json:"tags"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Version     uint           `gorm:"default:1" json:"-"`
}
