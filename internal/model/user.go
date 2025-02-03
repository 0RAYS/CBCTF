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

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"unique;not null" json:"name"`
	Password  string         `gorm:"not null" json:"-"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Country   string         `gorm:"default:'cn'" json:"country"`
	Avatar    string         `json:"-"`
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
		Contests int    `json:"contests"`
		Teams    int    `json:"teams"`
		Avatar   string `json:"avatar"`
	}{
		Tmp:      Tmp(m),
		Contests: len(m.Contests),
		Teams:    len(m.Teams),
		Avatar:   fmt.Sprintf("%s/%s", config.Env.Backend, m.Avatar),
	})
}

func InitUser(form constants.CreateUserForm) User {
	return User{
		Name:     form.Name,
		Password: utils.HashPassword(form.Password),
		Email:    form.Email,
		Country:  form.Country,
		Desc:     form.Desc,
		Avatar:   "",
		Verified: form.Verified,
		Hidden:   form.Hidden,
		Banned:   form.Banned,
	}
}
