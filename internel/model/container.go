package model

import (
	"gorm.io/gorm"
	"time"
)

type Container struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PodID     uint           `json:"pod_id"`
	Pod       Pod            `json:"-"`
	Name      string         `json:"name"`
	Image     string         `json:"image"`
	Hostname  string         `json:"hostname"`
	Flags     Strings        `gorm:"type:json" json:"flags"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Version   uint           `gorm:"default:1" json:"-"`
}
