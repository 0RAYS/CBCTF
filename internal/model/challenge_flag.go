package model

import (
	"CBCTF/internal/i18n"
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
	Challenge    Challenge      `json:"-"`
	DockerID     sql.Null[uint] `gorm:"default:null" json:"docker_id"`
	Docker       *Docker        `json:"-"`
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

func (c ChallengeFlag) CreateErrorString() string {
	return i18n.CreateChallengeFlagError
}

func (c ChallengeFlag) DeleteErrorString() string {
	return i18n.DeleteChallengeFlagError
}

func (c ChallengeFlag) GetErrorString() string {
	return i18n.GetChallengeFlagError
}

func (c ChallengeFlag) NotFoundErrorString() string {
	return i18n.ChallengeFlagNotFound
}

func (c ChallengeFlag) UpdateErrorString() string {
	return i18n.UpdateChallengeFlagError
}

func (c ChallengeFlag) GetUniqueKey() []string {
	return []string{"id"}
}

func (c ChallengeFlag) GetAllowedQueryFields() []string {
	return []string{}
}
