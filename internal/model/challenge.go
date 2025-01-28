package model

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"time"
)

var Static = 0
var Dynamic = 1
var Container = 2

var StaticFile = "attachment.zip"
var DynamicFile = "generator.zip"
var ContainerFile = "mounted.zip"

type Challenge struct {
	ID             string         `json:"id" gorm:"primaryKey"`
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
		ID:       utils.RandomString(),
		Name:     form.Name,
		Desc:     form.Desc,
		Flag:     form.Flag,
		Category: form.Category,
		Type:     form.Type,
		Path:     path,
	}
}

func (c *Challenge) GetFlag() string {
	switch c.Type {
	case Static:
		return c.Flag
	case Dynamic:
		return c.Flag
	case Container:
		return c.Flag
	default:
		return ""
	}
}
