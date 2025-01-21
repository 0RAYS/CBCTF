package model

import (
	"CBCTF/internal/model"
	"gorm.io/gorm"
	"time"
)

type DockerChallenge struct {
	ID         uint            `json:"id" gorm:"primaryKey"`
	Name       string          `json:"name" gorm:"unique;not null"`
	Msg        string          `json:"msg"`
	FlagPrefix string          `json:"flag_prefix"`
	Category   string          `json:"category"`
	Image      string          `json:"image"`
	Contests   []model.Contest `json:"-" gorm:"many2many:contest_challenges;"`
	CreatedAt  time.Time       `json:"-"`
	UpdatedAt  time.Time       `json:"-"`
	DeletedAt  gorm.DeletedAt  `json:"-" gorm:"index"`
}
