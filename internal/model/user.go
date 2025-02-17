package model

import (
	"CBCTF/internal/config"
	"CBCTF/internal/form"
	"CBCTF/internal/utils"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"index:idx_name_deleted,unique;not null" json:"name"`
	Password  string         `gorm:"not null" json:"-"`
	Email     string         `gorm:"index:idx_email_deleted,unique;not null" json:"email"`
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
	DeletedAt gorm.DeletedAt `gorm:"index;index:idx_name_deleted,unique;index:idx_email_deleted,unique" json:"-"`
}

func (m *User) MarshalJSON() ([]byte, error) {
	type Tmp User // 定义一个别名以避免递归调用
	return json.Marshal(&struct {
		Tmp
		Contests int    `json:"contests"`
		Teams    int    `json:"teams"`
		Avatar   string `json:"avatar"`
	}{
		Tmp:      Tmp(*m),
		Contests: len(m.Contests),
		Teams:    len(m.Teams),
		Avatar:   fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(m.Avatar, "/")),
	})
}

func InitUser(form form.CreateUserForm) User {
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
