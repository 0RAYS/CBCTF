package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
	ContestID  uint          `json:"contest_id"`
	Model      CheatRefModel `gorm:"default:null;type:json" json:"model"`
	Magic      string        `json:"magic"`
	IP         string        `json:"ip"`
	Reason     string        `json:"reason"`
	ReasonType string        `gorm:"index" json:"reason_type"`
	Type       string        `json:"type"`
	Checked    bool          `json:"checked"`
	Hash       string        `gorm:"type:varchar(32);index" json:"hash"`
	Comment    string        `json:"comment"`
	Time       time.Time     `json:"time"`
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
