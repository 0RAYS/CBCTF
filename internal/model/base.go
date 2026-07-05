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
	GetBaseModel() BaseModel
}

func (b BaseModel) GetBaseModel() BaseModel {
	return b
}

type RetVal struct {
	OK   bool
	Msg  string
	Attr map[string]any
	Data any
}

func SuccessRetVal(data ...any) RetVal {
	if len(data) > 0 {
		if len(data) == 1 {
			return RetVal{OK: true, Msg: i18n.Common.Success, Attr: nil, Data: data[0]}
		}
		return RetVal{OK: true, Msg: i18n.Common.Success, Attr: nil, Data: data}
	}
	return RetVal{OK: true, Msg: i18n.Common.Success, Attr: nil, Data: nil}
}
