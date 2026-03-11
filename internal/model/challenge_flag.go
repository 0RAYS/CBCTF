package model

import (
	"database/sql"
	"regexp"
)

type FlagInjectType string

var (
	StaticFlagTmpl  = regexp.MustCompile(`static\{(.*?)}`)
	UUIDFlagTmpl    = regexp.MustCompile(`uuid\{(.*?)}`)
	DynamicFlagTmpl = regexp.MustCompile(`dynamic\{(.*?)}`)

	EnvFlagInjectType    FlagInjectType = "env"
	VolumeFlagInjectType FlagInjectType = "volume"
)

// ChallengeFlag 题库中挑战的 flag 定义
// BelongsTo Challenge
// BelongsTo Docker
// HasMany ContestFlag
// HasMany TeamFlag
type ChallengeFlag struct {
	ChallengeID  uint           `json:"challenge_id"`
	Challenge    Challenge      `json:"-"`
	DockerID     sql.Null[uint] `gorm:"default:null" json:"docker_id"`
	Docker       Docker         `json:"-"`
	ContestFlags []ContestFlag  `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	TeamFlags    []TeamFlag     `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name         string         `json:"name"`
	Value        string         `json:"value"`
	InjectType   FlagInjectType `json:"inject_type"`
	Path         string         `json:"path"`
	BaseModel
}

func (c ChallengeFlag) TableName() string {
	return "challenge_flags"
}

func (c ChallengeFlag) ModelName() string {
	return "ChallengeFlag"
}

func (c ChallengeFlag) GetBaseModel() BaseModel {
	return c.BaseModel
}

func (c ChallengeFlag) UniqueFields() []string {
	return []string{"id"}
}

func (c ChallengeFlag) QueryFields() []string {
	return []string{}
}
