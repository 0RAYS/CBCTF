package model

import (
	"CBCTF/internal/config"
	"CBCTF/internal/form"
	"CBCTF/internal/utils"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
	"time"
)

const (
	Static    = "static"
	Dynamic   = "dynamic"
	Container = "container"

	AttachmentFile = "attachment.zip"
	GeneratorFile  = "generator.zip"
)

type Challenge struct {
	ID             string                 `json:"id" gorm:"primaryKey"`
	Name           string                 `json:"name" gorm:"not null"`
	Desc           string                 `json:"desc"`
	Flag           string                 `json:"flag"`
	Category       string                 `json:"category"`
	Type           string                 `json:"type" gorm:"default:'static'"`
	GeneratorImage string                 `json:"generator" gorm:"column:generator"`
	DockerImage    string                 `json:"docker" gorm:"column:docker"`
	Port           int32                  `json:"port" gorm:"default:8080"`
	CreatedAt      time.Time              `json:"-"`
	UpdatedAt      time.Time              `json:"-"`
	DeletedAt      gorm.DeletedAt         `json:"-" gorm:"index"`
	Version        optimisticlock.Version `json:"-" gorm:"default:1"`
}

// BasicDir 获取题目相关文件的目录
func (c *Challenge) BasicDir() string {
	return fmt.Sprintf("%s/challenges/%s", config.Env.Path, c.ID)
}

func InitChallenge(form form.CreateChallengeForm) Challenge {
	return Challenge{
		ID:             utils.UUID(),
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
