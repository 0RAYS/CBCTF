package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	Avatar  = "avatar"
	WriteUP = "writeup"
)

type File struct {
	ID        string         `gorm:"primarykey" json:"id"`
	Filename  string         `json:"filename"`
	Size      int64          `json:"size"`
	Path      string         `json:"-"`
	AdminID   uint           `json:"admin_id"`
	UserID    uint           `json:"user_id"`
	TeamID    uint           `json:"team_id"`
	ContestID uint           `json:"contest_id"`
	Suffix    string         `json:"suffix"`
	Hash      string         `json:"hash"`
	Type      string         `json:"type"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
