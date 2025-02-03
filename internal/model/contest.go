package model

import (
	"CBCTF/internal/config"
	"CBCTF/internal/constants"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Contest struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"unique;not null" json:"name"`
	Desc      string         `json:"desc"`
	Captcha   string         `json:"-"`
	Avatar    string         `json:"-"`
	Prefix    string         `json:"prefix" gorm:"default:'CBCTF'"`
	Size      int            `json:"size"`
	Start     time.Time      `json:"start"`
	Duration  time.Duration  `json:"duration"`
	Hidden    bool           `gorm:"default:false" json:"hidden"`
	Teams     []*Team        `json:"-"`
	Users     []*User        `gorm:"many2many:user_contests;" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c Contest) MarshalJSON() ([]byte, error) {
	type Tmp Contest // 定义一个别名以避免递归调用
	return json.Marshal(&struct {
		Tmp
		Users  int    `json:"users"`
		Teams  int    `json:"teams"`
		Avatar string `json:"avatar"`
	}{
		Tmp:    Tmp(c),
		Users:  len(c.Users),
		Teams:  len(c.Teams),
		Avatar: fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimSuffix(c.Avatar, "/")),
	})
}

func (c Contest) IsRunning() bool {
	return time.Now().After(c.Start) && time.Now().Before(c.Start.Add(c.Duration))
}

func InitContest(form constants.CreateContestForm) Contest {
	return Contest{
		Name:     form.Name,
		Desc:     form.Desc,
		Captcha:  form.Captcha,
		Avatar:   "",
		Size:     form.Size,
		Start:    form.Start,
		Hidden:   form.Hidden,
		Duration: form.Duration,
	}
}
