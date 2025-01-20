package model

import (
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"unique;not null" json:"name"`
	Password  string         `gorm:"not null" json:"-"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Website   string         `json:"website"`
	Country   string         `gorm:"default:'cn'" json:"country"`
	Avatar    string         `json:"avatar"`
	Desc      string         `json:"desc"`
	Verified  bool           `gorm:"default:false" json:"verified"`
	Hidden    bool           `gorm:"default:false" json:"hidden"`
	Banned    bool           `gorm:"default:false" json:"banned"`
	Teams     []*Team        `gorm:"many2many:user_teams;" json:"-"`
	Contests  []*Contest     `gorm:"many2many:user_contests;" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func InitUser(name string, password string, email string) User {
	return User{
		Name:     name,
		Password: utils.HashPassword(password),
		Email:    email,
		Website:  "",
		Country:  "cn",
		Desc:     "",
		Avatar:   "",
		Verified: false,
		Hidden:   false,
		Banned:   false,
	}
}
