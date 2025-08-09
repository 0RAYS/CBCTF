package repo

import (
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type EventRepo struct {
	BasicRepo[model.Event]
}

type CreateEventOptions struct {
	UserID             *uint
	TeamID             *uint
	ContestID          *uint
	ContestChallengeID *uint
	Desc               string
	Type               string
	IP                 string
	Magic              string
    Status             int
    Success            bool
    UserAgent          string
    TraceID            string
    Extra              string
}

func (c CreateEventOptions) Convert2Model() model.Model {
	return model.Event{
		UserID:             c.UserID,
		TeamID:             c.TeamID,
		ContestID:          c.ContestID,
		ContestChallengeID: c.ContestChallengeID,
		Desc:               c.Desc,
		Type:               c.Type,
		IP:                 c.IP,
		Magic:              c.Magic,
        Status:             c.Status,
        Success:            c.Success,
        UserAgent:          c.UserAgent,
        TraceID:            c.TraceID,
        Extra:              c.Extra,
	}
}

type UpdateEventOptions struct {
}

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
