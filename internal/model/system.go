package model

import (
	"gorm.io/gorm"
	"time"
)

type System struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	Key       string         `json:"key" gorm:"type:varchar(255);not null;unique_index"`
	Value     string         `json:"value" gorm:"type:text;not null"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
