package model

import (
	"CBCTF/internal/utils"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Team struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Desc      string         `json:"desc"`
	Captcha   string         `json:"-"`
	Avatar    string         `json:"avatar"`
	Banned    bool           `gorm:"default:false" json:"banned"`
	Hidden    bool           `gorm:"default:false" json:"hidden"`
	CaptainID uint           `json:"captain_id"`
	ContestID uint           `json:"contest_id"`
	Users     []*User        `gorm:"many2many:user_teams;" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m Team) MarshalJSON() ([]byte, error) {
	type Tmp Team // 定义一个别名以避免递归调用
	return json.Marshal(&struct {
		Tmp
		Users int `json:"users"`
	}{
		Tmp:   Tmp(m),
		Users: len(m.Users),
	})
}

func InitTeam(name string, captainID uint) Team {
	return Team{
		Name:      name,
		Desc:      "",
		Captcha:   utils.RandomString(),
		Avatar:    "",
		Banned:    false,
		Hidden:    false,
		CaptainID: captainID,
	}
}
