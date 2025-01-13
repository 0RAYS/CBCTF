package model

import (
	"RayWar/internal/utils"
	"gorm.io/gorm"
	"time"
)

type Contest struct {
	Name     string        `gorm:"unique;not null" json:"name"`
	Desc     string        `json:"desc"`
	Captcha  string        `json:"captcha"`
	Avatar   string        `json:"avatar"`
	Size     uint          `json:"size"`
	Start    time.Time     `json:"start"`
	Duration time.Duration `json:"duration"`
	Hidden   bool          `gorm:"default:true" json:"hidden"`
	Teams    []*Team       `gorm:"many2many:team_contests;"`
	gorm.Model
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
		Duration: 30 * 24 * time.Hour,
	}
}
