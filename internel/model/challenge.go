package model

import (
	"CBCTF/internel/config"
	"fmt"
)

const (
	StaticChallenge  = "static"
	DynamicChallenge = "dynamic"
	PodsChallenge    = "pods"

	AttachmentFile = "attachment.zip"
	GeneratorFile  = "generator.zip"
)

// Challenge 题目模型
// 题目的类型有三种: 静态题目, 动态题目, 容器题目
// 静态题目: flag 为 Flags 字段
// 动态题目: flag 为 Flags 字段
// 容器题目: flag 为 Dockers[].Flags 字段
type Challenge struct {
	ID          string       `gorm:"type:varchar(36);primaryKey" json:"id"`
	Usages      []Usage      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions []Submission `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name        string       `gorm:"not null" json:"name"`
	Desc        string       `json:"desc"`
	Category    string       `json:"category"`
	Type        string       `json:"type"`
	Generator   string       `json:"generator"`
	Flags       StringList   `gorm:"type:json" json:"flags"`
	Dockers     Dockers      `gorm:"type:json" json:"dockers"`
	BaseModel
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
		return fmt.Sprintf("%s/attachments/team-%d/%s.zip", config.Env.Path, teamID, c.ID)
	default:
		return c.StaticPath()
	}
}
