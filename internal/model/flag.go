package model

import (
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
	"time"
)

type Flag struct {
	ID          uint                   `json:"id" gorm:"primaryKey"`
	ContestID   uint                   `json:"contest_id"`
	TeamID      uint                   `json:"team_id"`
	ChallengeID string                 `json:"challenge_id"`
	Value       string                 `json:"value"`
	CreatedAt   time.Time              `json:"-"`
	UpdatedAt   time.Time              `json:"-"`
	DeletedAt   gorm.DeletedAt         `json:"-" gorm:"index"`
	Version     optimisticlock.Version `json:"-"`
}

func InitFlag(contestID, teamID uint, challengeID, value string) Flag {
	return Flag{
		ContestID:   contestID,
		TeamID:      teamID,
		ChallengeID: challengeID,
		Value:       value,
	}
}
