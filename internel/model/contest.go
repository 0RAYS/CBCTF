package model

import (
	"gorm.io/gorm"
	"time"
)

type Contest struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"index:idx_name_deleted,unique;not null" json:"name"`
	Desc        string         `json:"desc"`
	Captcha     string         `json:"-"`
	Avatar      string         `json:"avatar"`
	Prefix      string         `gorm:"default:'flag'" json:"prefix"`
	Size        int            `json:"size"`
	Start       time.Time      `json:"start"`
	Duration    time.Duration  `json:"-"`
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
	DeletedAt   gorm.DeletedAt `gorm:"index;index:idx_name_deleted,unique;" json:"-"`
	Version     uint           `gorm:"default:1" json:"-"`
}
