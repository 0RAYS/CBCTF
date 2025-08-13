package model

import (
	"time"

	"gorm.io/gorm"
)

type BasicModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Version   uint           `gorm:"default:1" json:"-"`
}

type Model interface {
	GetVersion() uint
	GetModelName() string
	GetBasicModel() BasicModel
	CreateErrorString() string
	DeleteErrorString() string
	GetErrorString() string
	NotFoundErrorString() string
	UpdateErrorString() string
	GetUniqueKey() []string
}
