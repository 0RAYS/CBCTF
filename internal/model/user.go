package model

import (
	"RayWar/internal/utils"
	"gorm.io/gorm"
)

type User struct {
	Name     string  `gorm:"unique;not null" json:"name"`
	Password string  `gorm:"not null" json:"password"`
	Email    string  `gorm:"unique;not null" json:"email"`
	Website  string  `json:"website"`
	Country  string  `gorm:"default:'cn'" json:"country"`
	Type     string  `gorm:"default:'user'" json:"type"`
	Avatar   string  `json:"avatar"`
	Desc     string  `gorm:"default:'不说话装高手'" json:"desc"`
	Verified bool    `gorm:"default:false" json:"verified"`
	Hidden   bool    `gorm:"default:false" json:"hidden"`
	Banned   bool    `gorm:"default:false" json:"banned"`
	Teams    []*Team `gorm:"many2many:user_teams;"`
	gorm.Model
}

func InitUser(name string, password string, email string) User {
	return User{
		Name:     name,
		Password: utils.HashPassword(password),
		Email:    email,
		Website:  "",
		Country:  "cn",
		Desc:     "",
		Type:     "user",
		Avatar:   "",
		Verified: false,
		Hidden:   false,
		Banned:   false,
	}
}
