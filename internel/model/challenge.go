package model

import (
	"CBCTF/internel/config"
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	StaticChallenge  = "static"
	DynamicChallenge = "dynamic"
	DockerChallenge  = "docker"
	DockersChallenge = "dockers"

	AttachmentFile = "attachment.zip"
	GeneratorFile  = "generator.zip"
)

type Challenge struct {
	ID          string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Desc        string         `json:"desc"`
	Category    string         `json:"category"`
	Type        string         `json:"type"`
	Generator   string         `json:"generator"`
	Flags       Strings        `gorm:"type:json" json:"flags"`
	Docker      Docker         `gorm:"type:json" json:"docker"`
	Dockers     Dockers        `gorm:"type:json" json:"dockers"`
	Usages      []Usage        `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions []Submission   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Version     uint           `gorm:"default:1" json:"-"`
}

// BasicDir 获取题目相关文件的目录
func (c *Challenge) BasicDir() string {
	return fmt.Sprintf("%s/challenges/%s", config.Env.Path, c.ID)
}

// StaticPath 获取静态题目文件的路径
func (c *Challenge) StaticPath() string {
	return fmt.Sprintf("%s/%s", c.BasicDir(), AttachmentFile)
}

// GeneratorPath 获取动态题目生成器的路径
func (c *Challenge) GeneratorPath() string {
	return fmt.Sprintf("%s/%s", c.BasicDir(), GeneratorFile)
}

// AttachmentPath 获取下载时, 题目附件的路径
func (c *Challenge) AttachmentPath(teamID uint) string {
	switch c.Type {
	case DynamicChallenge:
		return fmt.Sprintf("%s/attachments/%s/%d.zip", config.Env.Path, c.ID, teamID)
	default:
		return c.StaticPath()
	}
}
