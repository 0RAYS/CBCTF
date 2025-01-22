package model

import (
	"CBCTF/internal/utils"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Contest struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"unique;not null" json:"name"`
	Desc      string         `json:"desc"`
	Captcha   string         `json:"-"`
	Avatar    string         `json:"avatar"`
	Size      int            `json:"size"`
	Start     time.Time      `json:"start"`
	Duration  time.Duration  `json:"duration"`
	Hidden    bool           `gorm:"default:true" json:"hidden"`
	Teams     []*Team        `json:"-"`
	Users     []*User        `gorm:"many2many:user_contests;" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m Contest) MarshalJSON() ([]byte, error) {
	type Tmp Contest // 定义一个别名以避免递归调用
	return json.Marshal(&struct {
		Tmp
		Users int `json:"users"`
		Teams int `json:"teams"`
	}{
		Tmp:   Tmp(m),
		Users: len(m.Users),
		Teams: len(m.Teams),
	})
}

func InitContest(name string, desc string, size int, start time.Time, duration time.Duration, hidden bool) Contest {
	return Contest{
		Name:     name,
		Desc:     desc,
		Captcha:  utils.RandomString(),
		Avatar:   "",
		Size:     size,
		Start:    start,
		Hidden:   hidden,
		Duration: duration,
	}
}
