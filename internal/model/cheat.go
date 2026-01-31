package model

import (
	"database/sql"
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
)

type Cheat struct {
	UserID             sql.Null[uint] `gorm:"default:null" json:"user_id"`
	TeamID             sql.Null[uint] `gorm:"default:null" json:"team_id"`
	ContestID          sql.Null[uint] `gorm:"default:null" json:"contest_id"`
	ContestChallengeID sql.Null[uint] `gorm:"default:null" json:"contest_challenge_id"`
	ContestFlagID      sql.Null[uint] `gorm:"default:null" json:"contest_flag_id"`
	Magic              string         `json:"magic"`
	IP                 string         `json:"ip"`
	Reason             string         `json:"reason"`
	Type               string         `json:"type"`
	Checked            bool           `json:"checked"`
	Hash               string         `json:"hash"`
	Comment            string         `json:"comment"`
	Time               time.Time      `json:"time"`
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
	return []string{"id", "magic", "ip", "reason", "type", "checked", "hash", "comment", "time"}
}
