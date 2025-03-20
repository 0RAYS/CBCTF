package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Submission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UsageID     uint           `json:"usage_id"`
	ContestID   uint           `json:"contest_id"`
	ChallengeID string         `json:"challenge_id"`
	TeamID      uint           `json:"team_id"`
	UserID      uint           `json:"user_id"`
	Value       string         `json:"value"`
	Solved      bool           `json:"solved"`
	Score       float64        `json:"score" gorm:"default:0"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (s *Submission) MarshalJSON() ([]byte, error) {
	type Tmp Submission
	return json.Marshal(&struct {
		*Tmp
	}{
		Tmp: (*Tmp)(s),
	})
}

func InitSubmission(usage, contestID uint, challengeID string, teamID, userID uint, value string, solved bool, score float64) Submission {
	return Submission{
		UsageID:     usage,
		ContestID:   contestID,
		ChallengeID: challengeID,
		TeamID:      teamID,
		UserID:      userID,
		Value:       value,
		Solved:      solved,
		Score:       score,
	}
}
