package model

import (
	"CBCTF/internal/config"
	"CBCTF/internal/constants"
	"CBCTF/internal/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

var Static = 0
var Dynamic = 1
var Container = 2

var StaticFile = "attachment.zip"
var DynamicFile = "generator.zip"

type Challenge struct {
	ID             string         `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"not null"`
	Desc           string         `json:"desc"`
	Flag           string         `json:"flag"`
	Category       string         `json:"category"`
	Type           int            `json:"type" gorm:"default:0"`
	GeneratorImage string         `json:"generator" gorm:"column:generator"`
	DockerImage    string         `json:"docker" gorm:"column:docker"`
	Port           int32          `json:"port" gorm:"default:8080"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (c *Challenge) BasicDir() string {
	return fmt.Sprintf("%s/challenges/%s", config.Env.Gin.Upload.Path, c.ID)
}

func (c *Challenge) StaticPath() string {
	return fmt.Sprintf("%s/challenges/%s/%s", config.Env.Gin.Upload.Path, c.ID, StaticFile)
}

func (c *Challenge) GeneratorPath() string {
	return fmt.Sprintf("%s/challenges/%s/%s", config.Env.Gin.Upload.Path, c.ID, DynamicFile)
}

func (c *Challenge) AttachmentPath(teamID uint) string {
	switch c.Type {
	case Static:
		return c.StaticPath()
	case Dynamic:
		return fmt.Sprintf("%s/attachments/%s/%d.zip", config.Env.Gin.Upload.Path, c.ID, teamID)
	default:
		return c.StaticPath()
	}
}

func InitChallenge(form constants.CreateChallengeForm) Challenge {
	return Challenge{
		ID:             utils.RandomString(),
		Name:           form.Name,
		Desc:           form.Desc,
		Flag:           form.Flag,
		Category:       form.Category,
		Type:           form.Type,
		GeneratorImage: form.GeneratorImage,
		DockerImage:    form.DockerImage,
		Port:           form.Port,
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
