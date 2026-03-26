package model

import "regexp"

var (
	StaticFlagTmpl  = regexp.MustCompile(`static\{(.*?)}`)
	UUIDFlagTmpl    = regexp.MustCompile(`uuid\{(.*?)}`)
	DynamicFlagTmpl = regexp.MustCompile(`dynamic\{(.*?)}`)
)

// ChallengeFlag 题库中挑战的 flag 定义
// BelongsTo Challenge
// HasMany ContestFlag
// HasMany TeamFlag
type ChallengeFlag struct {
	ChallengeID  uint          `gorm:"index" json:"challenge_id"`
	Challenge    Challenge     `json:"-"`
	ContestFlags []ContestFlag `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	TeamFlags    []TeamFlag    `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name         string        `json:"name"`
	Value        string        `json:"value"`
	Binding      FlagBinding   `gorm:"type:jsonb" json:"binding"`
	BaseModel
}
