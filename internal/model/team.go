package model

import (
	"CBCTF/internal/utils"
	"gorm.io/gorm"
)

type Team struct {
	Name     string     `gorm:"not null" json:"name"`
	Desc     string     `gorm:"default:'不说话装酷'" json:"desc"`
	Captcha  string     `json:"captcha"`
	Avatar   string     `json:"avatar"`
	Banned   bool       `gorm:"default:false" json:"banned"`
	Hidden   bool       `gorm:"default:false" json:"hidden"`
	Users    []*User    `gorm:"many2many:user_teams;"`
	Contests []*Contest `gorm:"many2many:team_contests;"`
	gorm.Model
}

func InitTeam(name string) Team {
	return Team{
		Name:    name,
		Desc:    "",
		Captcha: utils.RandomString(),
		Avatar:  "",
		Banned:  false,
		Hidden:  false,
	}
}
