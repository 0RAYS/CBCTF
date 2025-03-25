package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	StaticChallenge  = "static"
	DynamicChallenge = "dynamic"
	DockerChallenge  = "docker"
	DockersChallenge = "dockers"
)

type Challenge struct {
	ID        string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Desc      string         `json:"desc"`
	Category  string         `json:"category"`
	Type      string         `json:"type"`
	Generator string         `json:"generator"`
	Flags     Strings        `gorm:"type:json" json:"flags"`
	Docker    Docker         `gorm:"type:json" json:"docker"`
	Dockers   Dockers        `gorm:"type:json" json:"dockers"`
	Usages    []Usage        `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Version   uint           `gorm:"default:1" json:"-"`
}
