package model

import (
	"time"
)

const (
	Cheater    = "cheater"
	Suspicious = "suspicious"
	Pass       = "pass"

	DifferentTokenMagic = "Device magic %s is different from token magic %s"
	SameDeviceMagic     = "%s has the same Device magic"
	ReqWebSameIP        = "%s request web with same IP"
	ReqVictimSameIP     = "%s request victim with same IP"
	SubmitOtherTeamFlag = "Team %d submitted flag of %s in Contest %d"

	ReasonTypeSameDevice   = "same_device"
	ReasonTypeSameWebIP    = "same_web_ip"
	ReasonTypeSameVictimIP = "same_victim_ip"
	ReasonTypeWrongFlag    = "wrong_flag"
	ReasonTypeTokenMagic   = "token_magic"
)

type Cheat struct {
	Model      UintMap   `gorm:"default:null;type:json" json:"model"`
	Magic      string    `json:"magic"`
	IP         string    `json:"ip"`
	Reason     string    `json:"reason"`
	ReasonType string    `gorm:"index" json:"reason_type"`
	Type       string    `json:"type"`
	Checked    bool      `json:"checked"`
	Hash       string    `gorm:"type:varchar(32);uniqueIndex" json:"hash"`
	Comment    string    `json:"comment"`
	Time       time.Time `json:"time"`
	BaseModel
}

func (c Cheat) ModelName() string {
	return "Cheat"
}

func (c Cheat) GetBaseModel() BaseModel {
	return c.BaseModel
}

func (c Cheat) UniqueFields() []string {
	return []string{"id", "hash"}
}

func (c Cheat) QueryFields() []string {
	return []string{"id", "magic", "ip", "reason", "reason_type", "type", "checked", "hash", "comment", "time"}
}
