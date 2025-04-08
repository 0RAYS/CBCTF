package model

import (
	"gorm.io/gorm"
	"time"
)

type Contest struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Desc        string         `json:"desc"`
	Captcha     string         `json:"captcha"`
	Avatar      string         `json:"avatar"`
	Prefix      string         `gorm:"default:'flag'" json:"prefix"`
	Size        int            `gorm:"default:4" json:"size"`
	Start       time.Time      `gorm:"not null" json:"start"`
	Duration    time.Duration  `json:"duration"`
	Blood       bool           `gorm:"default:true" json:"blood"`
	Hidden      bool           `gorm:"default:true" json:"hidden"`
	Rules       Strings        `gorm:"type:json" json:"rules"`
	Prizes      Prizes         `gorm:"type:json" json:"prizes"`
	Timelines   Timelines      `gorm:"type:json" json:"timelines"`
	Teams       []Team         `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Users       []*User        `gorm:"many2many:user_contests;" json:"-"`
	Notices     []Notice       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Usages      []Usage        `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Flags       []Flag         `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Cheats      []Cheat        `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions []Submission   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Version     uint           `gorm:"default:1" json:"-"`
}

func (c *Contest) IsOver() bool {
	return time.Now().After(c.Start.Add(c.Duration))
}

func (c *Contest) IsNotStart() bool {
	return time.Now().Before(c.Start)
}

func (c *Contest) IsRunning() bool {
	return (c.IsOver() || c.IsNotStart() || c.Hidden) != true
}

func (c *Contest) Status() string {
	if c.IsOver() {
		return "ContestIsOver"
	}
	if c.IsNotStart() {
		return "ContestNotRunning"
	}
	return "ContestIsRunning"
}
