package model

import (
	"time"
)

type CheatType string

type CheatReasonTmpl string

type CheatReasonType string

const (
	CheaterType    CheatType = "cheater"
	SuspiciousType CheatType = "suspicious"

	ReqWebSameIPTmpl        CheatReasonTmpl = "%s request web with same IP"
	ReqVictimSameIPTmpl     CheatReasonTmpl = "%s request victim with same IP"
	SubmitOtherTeamFlagTmpl CheatReasonTmpl = "Team %d submitted flag of %s in Contest %d"

	ReasonTypeSameWebIPType    CheatReasonType = "same_web_ip"
	ReasonTypeSameVictimIPType CheatReasonType = "same_victim_ip"
	ReasonTypeWrongFlagType    CheatReasonType = "wrong_flag"
)

type Cheat struct {
	Time       time.Time         `gorm:"default:null" json:"time"`
	Model      map[string][]uint `gorm:"serializer:json;type:jsonb;default:'{}'" json:"model"`
	IP         string            `json:"ip"`
	Reason     string            `json:"reason"`
	ReasonType CheatReasonType   `gorm:"index" json:"reason_type"`
	Type       CheatType         `json:"type"`
	Hash       string            `gorm:"type:varchar(32);uniqueIndex:idx_cheats_hash_active,where:deleted_at IS NULL;not null" json:"hash"`
	Comment    string            `json:"comment"`
	BaseModel
	ContestID uint `gorm:"index" json:"contest_id"`
	Checked   bool `gorm:"index" json:"checked"`
}
