package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Password    string         `gorm:"not null" json:"-"`
	Email       string         `gorm:"not null" json:"email"`
	Country     string         `gorm:"default:'CN'" json:"country"`
	Avatar      string         `json:"avatar"`
	Desc        string         `json:"desc"`
	Verified    bool           `gorm:"default:false" json:"verified"`
	Hidden      bool           `gorm:"default:false" json:"hidden"`
	Banned      bool           `gorm:"default:false" json:"banned"`
	Score       float64        `gorm:"default:0" json:"score"`
	Solved      int64          `gorm:"default:0" json:"solved"`
	Teams       []*Team        `gorm:"many2many:user_teams;" json:"-"`
	Contests    []*Contest     `gorm:"many2many:user_contests;" json:"-"`
	Submissions []Submission   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Containers  []Container    `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Devices     []Device       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Cheats      []Cheat        `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Version     uint           `gorm:"default:1" json:"-"`
}
