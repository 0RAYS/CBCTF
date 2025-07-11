package model

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"fmt"
)

const (
	StaticChallengeType   = "static"
	QuestionChallengeType = "question"
	DynamicChallengeType  = "dynamic"
	PodChallengeType      = "pod"
	VpcChallengeType      = "vpc"

	AttachmentFile = "attachment.zip"
	GeneratorFile  = "generator.zip"
)

// Challenge 题库中的挑战
// HasMany DockerGroup
// HasMany ChallengeFlag
// HasMany ContestChallenge
// HasMany Submission
type Challenge struct {
	DockerGroups      []DockerGroup      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ChallengeFlags    []ChallengeFlag    `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ContestChallenges []ContestChallenge `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions       []Submission       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	RandID            string             `gorm:"type:varchar(36);uniqueIndex;not null" json:"rand_id"`
	Name              string             `json:"name"`
	Desc              string             `json:"desc"`
	Category          string             `json:"category"`
	Type              string             `json:"type"`
	GeneratorImage    string             `json:"generator_image"`
	BasicModel
}

func (c Challenge) GetModelName() string {
	return "Challenge"
}

func (c Challenge) GetVersion() uint {
	return c.Version
}

func (c Challenge) CreateErrorString() string {
	return i18n.CreateChallengeError
}

func (c Challenge) DeleteErrorString() string {
	return i18n.DeleteChallengeError
}

func (c Challenge) GetErrorString() string {
	return i18n.GetChallengeError
}

func (c Challenge) NotFoundErrorString() string {
	return i18n.ChallengeNotFound
}

func (c Challenge) UpdateErrorString() string {
	return i18n.UpdateChallengeError
}

func (c Challenge) GetUniqueKey() []string {
	return []string{"id", "rand_id"}
}

func (c Challenge) GetForeignKeys() []string {
	return []string{"id"}
}

func (c Challenge) BasicDir() string {
	return fmt.Sprintf("%s/challenges/%d", config.Env.Path, c.ID)
}

// StaticPath 获取静态题目文件的路径
func (c Challenge) StaticPath() string {
	return fmt.Sprintf("%s/%s", c.BasicDir(), AttachmentFile)
}

// GeneratorPath 获取动态题目生成器的路径
func (c Challenge) GeneratorPath() string {
	return fmt.Sprintf("%s/%s", c.BasicDir(), GeneratorFile)
}

// AttachmentPath 获取下载时, 题目附件的路径
func (c Challenge) AttachmentPath(teamID uint) string {
	switch c.Type {
	case DynamicChallengeType:
		return fmt.Sprintf("%s/attachments/team-%d/%d.zip", config.Env.Path, teamID, c.ID)
	default:
		return c.StaticPath()
	}
}
