package db

import (
	"gorm.io/gorm"

	"CBCTF/internal/model"
)

type EventRepo struct {
	BaseRepo[model.Event]
}

func InitEventRepo(tx *gorm.DB) *EventRepo {
	return &EventRepo{
		DB: tx,
	}
}
