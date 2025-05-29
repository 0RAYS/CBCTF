package model

import (
	"time"
)

type Team struct {
	ContestID   uint         `gorm:"index:idx_name_contest,unique;not null" json:"contest_id"`
	Contest     Contest      `json:"-"`
	Users       []*User      `gorm:"many2many:user_teams;" json:"-"`
	Answers     []Answer     `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions []Submission `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Victims     []Victim     `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Cheats      []Cheat      `json:"-"`
	Name        string       `gorm:"type:varchar(255);index:idx_name_contest,unique;not null" json:"name"`
	Desc        string       `json:"desc"`
	Captcha     string       `json:"-"`
	Avatar      string       `json:"avatar"`
	Score       float64      `gorm:"default:0" json:"score"`
	Banned      bool         `gorm:"default:false" json:"banned"`
	Hidden      bool         `gorm:"default:false" json:"hidden"`
	CaptainID   uint         `json:"captain_id"`
	Rank        int          `gorm:"default:-1" json:"rank"`
	Last        time.Time    `gorm:"default:null" json:"last"`
	BaseModel
}
