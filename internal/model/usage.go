package model

import (
	"gorm.io/gorm"
	"time"
)

type Usage struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ContestID   uint           `json:"contest_id"`
	ChallengeID string         `json:"challenge_id"`
	Hidden      bool           `json:"hidden"`
	Score       int            `json:"score"`
	Flag        string         `json:"flag"`
	Attempt     int64          `json:"attempt" gorm:"default:0"`
	Hints       string         `json:"hints"`
	Tags        string         `json:"tags"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func InitUsage(challengeID string, contestID uint, flag string) Usage {
	return Usage{
		ContestID:   contestID,
		ChallengeID: challengeID,
		Flag:        flag,
	}
}
