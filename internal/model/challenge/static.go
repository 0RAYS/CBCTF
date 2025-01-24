package model

import (
	"CBCTF/internal/model"
	"gorm.io/gorm"
	"time"
)

type StaticChallenge struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	Name      string          `json:"name" gorm:"unique;not null"`
	Desc      string          `json:"desc"`
	Flag      string          `json:"flag"`
	Category  string          `json:"category"`
	Files     []model.Avatar  `json:"files" gorm:"many2many:challenge_files;"`
	Contests  []model.Contest `json:"-" gorm:"many2many:contest_challenges;"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
	DeletedAt gorm.DeletedAt  `json:"-" gorm:"index"`
}
