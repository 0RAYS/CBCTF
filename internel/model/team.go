package model

import (
	"gorm.io/gorm"
	"time"
)

type Team struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"index:idx_name_contest_id_deleted,unique;not null" json:"name"`
	ContestID   uint           `gorm:"index:idx_name_contest_id_deleted,unique;not null" json:"contest_id"`
	Contest     Contest        `json:"-"`
	Desc        string         `json:"desc"`
	Captcha     string         `json:"-"`
	Avatar      string         `json:"avatar"`
	Score       float64        `gorm:"default:0" json:"score"`
	Banned      bool           `gorm:"default:false" json:"banned"`
	Hidden      bool           `gorm:"default:false" json:"hidden"`
	CaptainID   uint           `json:"captain_id"`
	Rank        int            `gorm:"default:-1" json:"rank"`
	Last        time.Time      `json:"last"`
	Users       []*User        `gorm:"many2many:user_teams;" json:"-"`
	Answers     []Answer       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions []Submission   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Containers  []Container    `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Cheats      []Cheat        `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index;index:idx_name_contest_id_deleted,unique" json:"-"`
	Version     uint           `gorm:"default:1" json:"-"`
}
