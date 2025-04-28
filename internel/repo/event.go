package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type EventRepo struct {
	Repo[model.Event]
}

type CreateEventOptions struct {
	UserID    uint
	TeamID    uint
	ContestID uint
	UsageID   uint
	Desc      string
	Type      string
	IP        string
	Device    string
}

func InitEventRepo(tx *gorm.DB) *EventRepo {
	return &EventRepo{Repo: Repo[model.Event]{DB: tx, Model: "Event"}}
}

func (e *EventRepo) CountByKeyID(key string, id uint) (int64, bool, string) {
	var count int64
	res := e.DB.Model(&model.Event{}).Where(key+" = ?", id).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Events: %s", res.Error)
		return 0, false, ""
	}
	return count, true, "Success"
}

func (e *EventRepo) GetByKeyID(key string, id uint, limit, offset int, preloadL ...string) ([]model.Event, int64, bool, string) {
	var (
		events         = make([]model.Event, 0)
		count, ok, msg = e.CountByKeyID(key, id)
	)
	if !ok {
		return events, count, false, msg
	}
	res := e.DB.Model(&model.Event{}).Where(key+" = ?", id)
	res = preload(res, preloadL...).Limit(limit).Offset(offset).Find(&events)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Events: %s", res.Error)
		return events, count, false, "GetEventError"
	}
	return events, count, true, "Success"
}
