package model

import (
	"CBCTF/internal/i18n"
	"time"
)

// Team
// BelongsTo Contest
// ManyToMany User
// ManyToMany Submission
// HasMany TeamFlag
type Team struct {
	ContestID   uint         `gorm:"index:idx_name_contest,unique;not null" json:"contest_id"`
	Contest     Contest      `json:"-"`
	Users       []*User      `gorm:"many2many:user_teams;" json:"-"`
	Submissions []Submission `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	TeamFlags   []TeamFlag   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name        string       `gorm:"type:varchar(255);index:idx_name_contest,unique;not null" json:"name"`
	Desc        string       `json:"desc"`
	Captcha     string       `json:"-"`
	Avatar      AvatarURL    `json:"avatar"`
	Score       float64      `gorm:"default:0" json:"score"`
	Banned      bool         `gorm:"default:false" json:"banned"`
	Hidden      bool         `gorm:"default:false" json:"hidden"`
	CaptainID   uint         `json:"captain_id"`
	Rank        int          `gorm:"default:-1" json:"rank"`
	Last        time.Time    `gorm:"default:null" json:"last"`
	UserCount   int64        `json:"user_count"`
	BasicModel
}

func (t Team) GetModelName() string {
	return "Team"
}

func (t Team) GetBasicModel() BasicModel {
	return t.BasicModel
}

func (t Team) CreateErrorString() string {
	return i18n.CreateTeamError
}

func (t Team) DeleteErrorString() string {
	return i18n.DeleteTeamError
}

func (t Team) GetErrorString() string {
	return i18n.GetTeamError
}

func (t Team) NotFoundErrorString() string {
	return i18n.TeamNotFound
}

func (t Team) UpdateErrorString() string {
	return i18n.UpdateTeamError
}

func (t Team) GetUniqueKey() []string {
	return []string{"id"}
}
