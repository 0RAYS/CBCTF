package model

import (
	"CBCTF/internal/config"
	"CBCTF/internal/constants"
	"CBCTF/internal/utils"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Team struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Desc      string         `json:"desc"`
	Captcha   string         `json:"-"`
	Avatar    string         `json:"-"`
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
		Users  int    `json:"users"`
		Avatar string `json:"avatar"`
	}{
		Tmp:    Tmp(m),
		Users:  len(m.Users),
		Avatar: fmt.Sprintf("%s/%s", config.Env.Backend, m.Avatar),
	})
}

func InitTeam(form constants.CreateTeamForm, captainID uint) Team {
	captcha := utils.RandomString()
	if form.Captcha != "" {
		captcha = form.Captcha
	}
	return Team{
		Name:      form.Name,
		Desc:      form.Desc,
		Captcha:   captcha,
		Avatar:    "",
		Banned:    false,
		Hidden:    false,
		CaptainID: captainID,
	}
}
