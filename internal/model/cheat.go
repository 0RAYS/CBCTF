package model

import (
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"time"
)

const (
	Suspect string = "suspect"
	Cheater string = "cheater"
)

const (
	MagicNotMatch string = "Magic in token and request headers are not matched"
	SameFlag      string = "Flag is same"
)

type Cheat struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id"`
	TeamID       uint           `json:"team_id"`
	ContestID    uint           `json:"contest_id"`
	Reason       string         `json:"reason"`
	Type         string         `json:"type"`
	Associations utils.Strings  `json:"associations"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func InitCheat(userID, teamID, contestID uint, reason string, t string) Cheat {
	return Cheat{
		ID:           utils.UUID(),
		UserID:       userID,
		TeamID:       teamID,
		ContestID:    contestID,
		Reason:       reason,
		Type:         t,
		Associations: []string{},
	}
}
