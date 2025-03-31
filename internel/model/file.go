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
	Uploader  uint           `json:"uploader"`
	Suffix    string         `json:"suffix"`
	Hash      string         `json:"hash"`
	Type      string         `json:"type"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
