package model

import (
	"CBCTF/internal/i18n"
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
	UserID             sql.Null[uint]    `gorm:"default:null" json:"user_id"`
	User               *User             `json:"-"`
	TeamID             sql.Null[uint]    `gorm:"default:null" json:"team_id"`
	Team               *Team             `json:"-"`
	ContestID          sql.Null[uint]    `gorm:"default:null" json:"contest_id"`
	Contest            *Contest          `json:"-"`
	ContestChallengeID sql.Null[uint]    `gorm:"default:null" json:"contest_challenge_id"`
	ContestChallenge   *ContestChallenge `json:"-"`
	ContestFlagID      sql.Null[uint]    `gorm:"default:null" json:"contest_flag_id"`
	ContestFlag        *ContestFlag      `json:"-"`
	Magic              string            `json:"magic"`
	IP                 string            `json:"ip"`
	Reason             string            `json:"reason"`
	Type               string            `json:"type"`
	Checked            bool              `json:"checked"`
	Hash               string            `json:"hash"`
	Comment            string            `json:"comment"`
	Time               time.Time         `json:"time"`
	BaseModel
}

func (c Cheat) GetModelName() string {
	return "Cheat"
}

func (c Cheat) GetBaseModel() BaseModel {
	return c.BaseModel
}

func (c Cheat) CreateErrorString() string {
	return i18n.CreateCheatError
}

func (c Cheat) DeleteErrorString() string {
	return i18n.DeleteCheatError
}

func (c Cheat) GetErrorString() string {
	return i18n.GetCheatError
}

func (c Cheat) NotFoundErrorString() string {
	return i18n.CheatNotFound
}

func (c Cheat) UpdateErrorString() string {
	return i18n.UpdateCheatError
}

func (c Cheat) GetUniqueKey() []string {
	return []string{"id", "hash"}
}

func (c Cheat) GetAllowedQueryFields() []string {
	return []string{"id", "magic", "ip", "reason", "type", "checked", "hash", "comment", "time"}
}
