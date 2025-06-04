package repo

import (
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type EventRepo struct {
	Basic[model.Event]
}

type CreateEventOptions struct {
	Desc      string
	Type      string
	IP        string
	Magic     string
	Reference model.Reference
}

func (c CreateEventOptions) Convert2Model() model.Model {
	return model.Event{
		Desc:      c.Desc,
		Type:      c.Type,
		IP:        c.IP,
		Magic:     c.Magic,
		Reference: c.Reference,
	}
}

type UpdateEventOptions struct {
}

func (u UpdateEventOptions) Convert2Map() map[string]any {
	return make(map[string]any)
}

func InitEventRepo(tx *gorm.DB) *EventRepo {
	return &EventRepo{
		Basic: Basic[model.Event]{
			DB: tx,
		},
	}
}
