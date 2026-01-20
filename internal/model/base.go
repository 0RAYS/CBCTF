package model

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type BaseModel struct {
	ID        uint                   `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time              `json:"-"`
	UpdatedAt time.Time              `json:"-"`
	DeletedAt gorm.DeletedAt         `gorm:"index" json:"-"`
	Version   optimisticlock.Version `json:"-"`
}

type Model interface {
	GetModelName() string
	GetBaseModel() BaseModel
	CreateErrorString() string
	DeleteErrorString() string
	GetErrorString() string
	NotFoundErrorString() string
	UpdateErrorString() string
	GetUniqueKey() []string
}
