package model

import (
	"gorm.io/gorm"
	"time"
)

type StaticChallenge struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"unique;not null"`
	Desc      string         `json:"desc"`
	Flag      string         `json:"flag"`
	Category  string         `json:"category"`
	Path      string         `json:"path"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
