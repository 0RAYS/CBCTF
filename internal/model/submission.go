package model

import (
	"gorm.io/gorm"
	"time"
)

type Submission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ContestID   uint           `json:"contest_id"`
	ChallengeID string         `json:"challenge_id"`
	TeamID      uint           `json:"team_id"`
	UserID      uint           `json:"user_id"`
	Value       string         `json:"value"`
	Solved      bool           `json:"solved"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func InitSubmission(contestID uint, challengeID string, teamID, userID uint, value string, solved bool) Submission {
	return Submission{
		ContestID:   contestID,
		ChallengeID: challengeID,
		TeamID:      teamID,
		UserID:      userID,
		Value:       value,
		Solved:      solved,
	}
}
