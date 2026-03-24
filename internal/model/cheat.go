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
	ContestID  uint            `gorm:"index" json:"contest_id"`
	Model      CheatRefModel   `gorm:"default:null;type:jsonb" json:"model"`
	Magic      string          `json:"magic"`
	IP         string          `json:"ip"`
	Reason     string          `json:"reason"`
	ReasonType CheatReasonType `gorm:"index" json:"reason_type"`
	Type       CheatType       `json:"type"`
	Checked    bool            `gorm:"index" json:"checked"`
	Hash       string          `gorm:"type:varchar(32);index" json:"hash"`
	Comment    string          `json:"comment"`
	Time       time.Time       `gorm:"default:null" json:"time"`
	BaseModel
}

type CheatRefModel map[string][]uint

func (c CheatRefModel) Value() (driver.Value, error) {
	if len(c) == 0 {
		return nil, nil
	}
	return json.Marshal(c)
}

func (c *CheatRefModel) Scan(value any) error {
	if err := scanJSON(value, c); err != nil {
		return fmt.Errorf("failed to scan CheatRefModel value")
	}
	return nil
}
