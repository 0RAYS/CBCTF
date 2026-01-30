package model

import (
	"database/sql"
	"regexp"
)

var (
	StaticFlag  = regexp.MustCompile(`static\{(.*?)}`)
	UUIDFlag    = regexp.MustCompile(`uuid\{(.*?)}`)
	DynamicFlag = regexp.MustCompile(`dynamic\{(.*?)}`)

	EnvInjectType    = "env"
	VolumeInjectType = "volume"
)

// ChallengeFlag 题库中挑战的 flag 定义
// BelongsTo Challenge
// BelongsTo Docker
// HasMany ContestFlag
type ChallengeFlag struct {
	ChallengeID  uint           `json:"challenge_id"`
	DockerID     sql.Null[uint] `gorm:"default:null" json:"docker_id"`
	ContestFlags []ContestFlag  `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	TeamFlags    []TeamFlag     `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name         string         `json:"name"`
	Value        string         `json:"value"`
	InjectType   string         `json:"inject_type"`
	Path         string         `json:"path"`
	BaseModel
}

func (c ChallengeFlag) GetModelName() string {
	return "ChallengeFlag"
}

func (c ChallengeFlag) GetBaseModel() BaseModel {
	return c.BaseModel
}

func (c ChallengeFlag) GetUniqueKey() []string {
	return []string{"id"}
}

func (c ChallengeFlag) GetAllowedQueryFields() []string {
	return []string{}
}
