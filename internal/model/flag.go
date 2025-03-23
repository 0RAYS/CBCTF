package model

import (
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
	"time"
)

type Flag struct {
	ID          uint                   `json:"id" gorm:"primaryKey"`
	ContestID   uint                   `json:"contest_id"`
	TeamID      uint                   `json:"team_id"`
	ChallengeID string                 `json:"challenge_id"`
	Values      utils.Strings          `json:"values" gorm:"type:json"`
	CreatedAt   time.Time              `json:"-"`
	UpdatedAt   time.Time              `json:"-"`
	DeletedAt   gorm.DeletedAt         `json:"-" gorm:"index"`
	Version     optimisticlock.Version `json:"-" gorm:"default:1"`
}

func InitFlag(contestID, teamID uint, challengeID string, values ...string) Flag {
	return Flag{
		ContestID:   contestID,
		TeamID:      teamID,
		ChallengeID: challengeID,
		Values:      values,
	}
}
