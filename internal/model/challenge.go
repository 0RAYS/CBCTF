package model

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"time"
)

type Challenge struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Desc      string         `json:"desc"`
	Flag      string         `json:"flag"`
	Category  string         `json:"category"`
	Path      string         `json:"path"`
	Type      string         `json:"type"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func InitChallenge(form constants.CreateChallengeForm) Challenge {
	return Challenge{
		Name:     form.Name,
		Desc:     form.Desc,
		Flag:     form.Flag,
		Category: form.Category,
		Type:     form.Type,
		Path:     utils.RandomString(),
	}
}
