package model

import (
	"CBCTF/internal/utils"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"unique;not null" json:"name"`
	Password  string         `gorm:"not null" json:"-"`
	Email     string         `gorm:"unique;not null" json:"email"`
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

func (m User) MarshalJSON() ([]byte, error) {
	type Tmp User // 定义一个别名以避免递归调用
	return json.Marshal(&struct {
		Tmp
		Contests int `json:"contests"`
		Teams    int `json:"teams"`
	}{
		Tmp:      Tmp(m),
		Contests: len(m.Contests),
		Teams:    len(m.Teams),
	})
}

func InitUser(name string, password string, email string, desc string, country string, hidden bool, verified bool, banned bool) User {
	return User{
		Name:     name,
		Password: utils.HashPassword(password),
		Email:    email,
		Country:  country,
		Desc:     desc,
		Avatar:   "",
		Verified: verified,
		Hidden:   hidden,
		Banned:   banned,
	}
}
