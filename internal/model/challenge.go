package model

import (
	"CBCTF/internal/constants"
	"gorm.io/gorm"
	"time"
)

var Static int = 0
var Dynamic int = 1
var Container int = 2

type Challenge struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"not null"`
	Desc           string         `json:"desc"`
	Flag           string         `json:"flag"`
	Category       string         `json:"category"`
	Path           string         `json:"path"`
	Type           int            `json:"type" gorm:"default:0"`
	GeneratorImage string         `json:"generator" gorm:"column:generator"`
	DockerImage    string         `json:"docker" gorm:"column:docker"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func InitChallenge(form constants.CreateChallengeForm, path string) Challenge {
	return Challenge{
		Name:     form.Name,
		Desc:     form.Desc,
		Flag:     form.Flag,
		Category: form.Category,
		Type:     form.Type,
		Path:     path,
	}
}
