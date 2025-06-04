package model

import (
	"gorm.io/gorm"
	"time"
)

type Model interface {
	GetID() uint
	GetVersion() uint
	GetModelName() string
	CreateErrorString() string
	DeleteErrorString() string
	GetErrorString() string
	NotFoundErrorString() string
	UpdateErrorString() string
	GetUniqueKey() []string
}

type Basic struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Version   uint           `gorm:"default:1" json:"-"`
}
