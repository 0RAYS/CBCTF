package model

import (
	"CBCTF/internal/utils"
	"gorm.io/gorm"
)

type Team struct {
	ID        uint    `gorm:"primarykey"`
	Name      string  `gorm:"not null" json:"name"`
	Desc      string  `json:"desc"`
	Password  string  `json:"-"`
	Avatar    string  `json:"avatar"`
	Banned    bool    `gorm:"default:false" json:"banned"`
	Hidden    bool    `gorm:"default:false" json:"hidden"`
	CaptainID uint    `json:"captain_id"`
	ContestID uint    `json:"contest_id"`
	Users     []*User `gorm:"many2many:user_teams;" json:"-"`
	gorm.Model
}

func InitTeam(name string, captainID uint) Team {
	return Team{
		Name:      name,
		Desc:      "",
		Password:  utils.RandomString(),
		Avatar:    "",
		Banned:    false,
		Hidden:    false,
		CaptainID: captainID,
	}
}
