package db

import (
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type EventRepo struct {
	BaseRepo[model.Event]
}

type CreateEventOptions struct {
	IsAdmin bool
	Type    string
	Success bool
	IP      string
	Magic   string
	Models  model.UintMap
}

func (c CreateEventOptions) Convert2Model() model.Model {
	return model.Event{
		IsAdmin: c.IsAdmin,
		Type:    c.Type,
		Success: c.Success,
		IP:      c.IP,
		Magic:   c.Magic,
		Models:  c.Models,
	}
}

func InitEventRepo(tx *gorm.DB) *EventRepo {
	return &EventRepo{
		BaseRepo: BaseRepo[model.Event]{
			DB: tx,
		},
	}
}
