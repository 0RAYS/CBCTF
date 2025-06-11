package model

import "CBCTF/internel/i18n"

const (
	Cheater    = "cheater"
	Suspicious = "suspicious"

	DifferentTokenMagic = "Device magic %s is different from token magic %s"
	SameDeviceMagic     = "User %d has the same device magic as user %d"
	SameRequestIP       = "%s has the same request IP %s"
)

type Cheat struct {
	UserID             *uint             `gorm:"default:null" json:"user_id"`
	User               *User             `json:"-"`
	TeamID             *uint             `gorm:"default:null" json:"team_id"`
	Team               *Team             `json:"-"`
	ContestID          *uint             `gorm:"default:null" json:"contest_id"`
	Contest            *Contest          `json:"-"`
	ContestChallengeID *uint             `gorm:"default:null" json:"contest_challenge_id"`
	ContestChallenge   *ContestChallenge `json:"-"`
	ContestFlagID      *uint             `gorm:"default:null" json:"contest_flag_id"`
	ContestFlag        *ContestFlag      `json:"-"`
	Magic              string            `json:"magic"`
	IP                 string            `json:"ip"`
	Reason             string            `json:"reason"`
	Type               string            `json:"type"`
	Checked            bool              `json:"checked"`
	Hash               string            `gorm:"type:varchar(32);uniqueIndex;not null" json:"hash"`
	References         UintList          `gorm:"type:json" json:"references"`
	Basic
}

func (c Cheat) GetModelName() string {
	return "Cheat"
}

func (c Cheat) GetID() uint {
	return c.ID
}

func (c Cheat) GetVersion() uint {
	return c.Version
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
	return []string{"id"}
}
