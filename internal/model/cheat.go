package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type CheatType string

type CheatReasonTmpl string

type CheatReasonType string

const (
	CheaterType    CheatType = "cheater"
	SuspiciousType CheatType = "suspicious"
	PassType       CheatType = "pass"

	DifferentTokenMagicTmpl CheatReasonTmpl = "Device magic %s is different from token magic %s"
	SameDeviceMagicTmpl     CheatReasonTmpl = "%s has the same Device magic"
	ReqWebSameIPTmpl        CheatReasonTmpl = "%s request web with same IP"
	ReqVictimSameIPTmpl     CheatReasonTmpl = "%s request victim with same IP"
	SubmitOtherTeamFlagTmpl CheatReasonTmpl = "Team %d submitted flag of %s in Contest %d"

	ReasonTypeSameDeviceType   CheatReasonType = "same_device"
	ReasonTypeSameWebIPType    CheatReasonType = "same_web_ip"
	ReasonTypeSameVictimIPType CheatReasonType = "same_victim_ip"
	ReasonTypeWrongFlagType    CheatReasonType = "wrong_flag"
	ReasonTypeTokenMagicType   CheatReasonType = "token_magic"
)

type Cheat struct {
	ContestID  uint            `json:"contest_id"`
	Model      CheatRefModel   `gorm:"default:null;type:json" json:"model"`
	Magic      string          `json:"magic"`
	IP         string          `json:"ip"`
	Reason     string          `json:"reason"`
	ReasonType CheatReasonType `gorm:"index" json:"reason_type"`
	Type       CheatType       `json:"type"`
	Checked    bool            `json:"checked"`
	Hash       string          `gorm:"type:varchar(32);index" json:"hash"`
	Comment    string          `json:"comment"`
	Time       time.Time       `json:"time"`
	BaseModel
}

func (c Cheat) TableName() string {
	return "cheats"
}

func (c Cheat) ModelName() string {
	return "Cheat"
}

func (c Cheat) GetBaseModel() BaseModel {
	return c.BaseModel
}

func (c Cheat) UniqueFields() []string {
	return []string{"id"}
}

func (c Cheat) QueryFields() []string {
	return []string{"id", "magic", "ip", "reason", "reason_type", "type", "checked", "hash", "comment", "time", "contest_id"}
}

type CheatRefModel map[string][]uint

func (c CheatRefModel) Value() (driver.Value, error) {
	if len(c) == 0 {
		return nil, nil
	}
	return json.Marshal(c)
}

func (c *CheatRefModel) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan CheatRefModel value")
	}
	if len(bytes) == 0 {
		*c = nil
		return nil
	}
	return json.Unmarshal(bytes, c)
}
