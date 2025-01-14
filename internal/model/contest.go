package model

import (
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"time"
)

type Contest struct {
	ID       uint          `gorm:"primarykey"`
	Name     string        `gorm:"unique;not null" json:"name"`
	Desc     string        `json:"desc"`
	Password string        `json:"password"`
	Avatar   string        `json:"avatar"`
	Size     int           `json:"size"`
	Start    time.Time     `json:"start"`
	Duration time.Duration `json:"duration"`
	Hidden   bool          `gorm:"default:true" json:"hidden"`
	Teams    []*Team       `json:"-"`
	Users    []*User       `gorm:"many2many:user_contests;" json:"-"`
	gorm.Model
}

func InitContest(name string) Contest {
	return Contest{
		Name:     name,
		Desc:     "",
		Password: utils.RandomString(),
		Avatar:   "",
		Size:     1,
		Start:    time.Now(),
		Hidden:   true,
		Duration: 2 * 24 * time.Hour,
	}
}
