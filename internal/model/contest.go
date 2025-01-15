package model

import (
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"time"
)

type Contest struct {
	ID        uint           `gorm:"primarykey"`
	Name      string         `gorm:"unique;not null" json:"name"`
	Desc      string         `json:"desc"`
	Captcha   string         `json:"captcha"`
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

func InitContest(name string) Contest {
	return Contest{
		Name:     name,
		Desc:     "",
		Captcha:  utils.RandomString(),
		Avatar:   "",
		Size:     1,
		Start:    time.Now(),
		Hidden:   true,
		Duration: 2 * 24 * time.Hour,
	}
}
