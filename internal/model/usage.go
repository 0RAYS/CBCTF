package model

import (
	"CBCTF/internal/constants"
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
	Attempt     int            `json:"attempt" gorm:"default:0"`
	Hint        string         `json:"hint"`
	Tag         string         `json:"tag"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func InitUsage(form constants.CreateUsageForm, contestID uint, flag string) Usage {
	return Usage{
		ContestID:   contestID,
		ChallengeID: form.ChallengeID,
		Hidden:      form.Hidden,
		Score:       form.Score,
		Flag:        flag,
		Attempt:     form.Attempt,
		Hint:        form.Hint,
		Tag:         form.Tag,
	}
}
