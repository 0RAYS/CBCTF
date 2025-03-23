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
	Static  = "static"
	Dynamic = "dynamic"
	Docker  = "docker"
	Dockers = "dockers"

	AttachmentFile = "attachment.zip"
	GeneratorFile  = "generator.zip"
)

type Challenge struct {
	ID             string                 `json:"id" gorm:"primaryKey"`
	Name           string                 `json:"name" gorm:"not null"`
	Desc           string                 `json:"desc"`
	Flag           string                 `json:"flag"`
	Flags          utils.Strings          `json:"flags" gorm:"type:json"`
	Category       string                 `json:"category"`
	Type           string                 `json:"type" gorm:"default:'static'"`
	GeneratorImage string                 `json:"generator" gorm:"column:generator"`
	DockerImage    string                 `json:"docker" gorm:"column:docker"`
	Port           int32                  `json:"port" gorm:"default:8080"`
	Dockers        utils.Dockers          `json:"dockers" gorm:"type:json"`
	CreatedAt      time.Time              `json:"-"`
	UpdatedAt      time.Time              `json:"-"`
	DeletedAt      gorm.DeletedAt         `json:"-" gorm:"index"`
	Version        optimisticlock.Version `json:"-" gorm:"default:1"`
}

// BasicDir 获取题目相关文件的目录
func (c Challenge) BasicDir() string {
	return fmt.Sprintf("%s/challenges/%s", config.Env.Path, c.ID)
}

func InitChallenge(form form.CreateChallengeForm) Challenge {
	c := Challenge{
		ID:       utils.UUID(),
		Name:     form.Name,
		Desc:     form.Desc,
		Flag:     form.Flag,
		Flags:    utils.Strings{form.Flag},
		Category: form.Category,
		Type:     form.Type,
	}
	switch form.Type {
	case Static:
		return c
	case Dynamic:
		c.GeneratorImage = form.GeneratorImage
		return c
	case Docker, Dockers:
		c.DockerImage = form.DockerImage
		c.Port = form.Port
		c.Dockers = utils.Dockers{{Image: form.DockerImage, Ports: []int32{form.Port}}}
		return c
	default:
		return c
	}
}
