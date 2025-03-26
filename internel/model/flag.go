package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	StaticScore      uint = 0
	LinearScore      uint = 1
	LogarithmicScore uint = 2
)

type Flag struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ContestID    uint           `json:"contest_id"`
	UsageID      string         `json:"usage_id"`
	Value        string         `json:"value"`
	Score        float64        `gorm:"default:1000" json:"score"`
	CurrentScore float64        `gorm:"default:1000" json:"current_score"`
	Decay        float64        `gorm:"default:50" json:"decay"`
	MinScore     float64        `gorm:"default:100" json:"min_score"`
	ScoreType    uint           `gorm:"default:0" json:"score_type"`
	Solvers      int64          `json:"solvers"`
	Attempt      int64          `json:"attempt"`
	Blood        Uints          `gorm:"type:json" json:"blood"`
	Answers      []Answer       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions  []Submission   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Last         time.Time      `json:"last"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Version      uint           `gorm:"default:1" json:"-"`
}
