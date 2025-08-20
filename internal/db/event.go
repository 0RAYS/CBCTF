package db

import (
	"CBCTF/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventRepo struct {
	BasicRepo[model.Event]
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

type UpdateEventOptions struct{}

func (u UpdateEventOptions) Convert2Map() map[string]any {
	return make(map[string]any)
}

func InitEventRepo(tx *gorm.DB) *EventRepo {
	return &EventRepo{
		BasicRepo: BasicRepo[model.Event]{
			DB: tx,
		},
	}
}
