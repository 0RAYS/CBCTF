package db

import (
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type EventRepo struct {
	BaseRepo[model.Event]
}

func InitEventRepo(tx *gorm.DB) *EventRepo {
	return &EventRepo{
		BaseRepo: BaseRepo[model.Event]{
			DB: tx,
		},
	}
}
