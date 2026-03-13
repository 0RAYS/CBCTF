package model

import (
	"database/sql"
	"time"
)

// Generator
// BelongsTo Challenge
// BelongsTo Contest
type Generator struct {
	ChallengeID uint           `json:"challenge_id"`
	Challenge   Challenge      `json:"-"`
	ContestID   sql.Null[uint] `json:"contest_id"`
	Contest     Contest        `json:"-"`
	Name        string         `json:"pod_name"`
	Success     int64          `json:"success"`
	SuccessLast time.Time      `gorm:"default:null" json:"success_last"`
	Failure     int64          `json:"failure"`
	FailureLast time.Time      `gorm:"default:null" json:"failure_last"`
	BaseModel
}

func (g Generator) TableName() string {
	return "generators"
}

func (g Generator) ModelName() string {
	return "Generator"
}

func (g Generator) GetBaseModel() BaseModel {
	return g.BaseModel
}

func (g Generator) UniqueFields() []string {
	return []string{"id"}
}

func (g Generator) QueryFields() []string {
	return []string{"id", "challenge_id", "contest_id", "success_count", "failure_count", "period"}
}
