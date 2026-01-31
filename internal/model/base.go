package model

import (
	"CBCTF/internal/i18n"
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
	ModelName() string
	GetBaseModel() BaseModel
	UniqueFields() []string
	QueryFields() []string
}

type RetVal struct {
	OK   bool
	Msg  string
	Attr map[string]any
	Data any
}

func SuccessRetVal(data ...any) RetVal {
	if len(data) > 0 {
		return RetVal{true, i18n.Common.Success, nil, data[0]}
	}
	return RetVal{true, i18n.Common.Success, nil, nil}
}
