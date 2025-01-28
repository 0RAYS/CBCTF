package model

import (
	"gorm.io/gorm"
	"time"
)

type Flag struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ContestID   uint           `json:"contest_id"`
	TeamID      uint           `json:"team_id"`
	ChallengeID string         `json:"challenge_id"`
	Value       string         `json:"value"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func InitFlag(contestID, teamID uint, challengeID, value string) Flag {
	return Flag{
		ContestID:   contestID,
		TeamID:      teamID,
		ChallengeID: challengeID,
		Value:       value,
	}
}
