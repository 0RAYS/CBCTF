package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"gorm.io/gorm"
	"time"
)

type TrafficRepo struct {
	Repo[model.Traffic]
}

type CreateTrafficOptions struct {
	VictimID uint
	PodID    uint
	SrcIP    string
	DstIP    string
	SrcPort  uint16
	DstPort  uint16
	Payload  string
	Time     time.Time
	Type     string
	Path     string
}

func InitTrafficRepo(tx *gorm.DB) *TrafficRepo {
	return &TrafficRepo{Repo: Repo[model.Traffic]{DB: tx, Model: "Traffic"}}
}

func (t *TrafficRepo) CountByKey(key string, id uint) (int64, bool, string) {
	var count int64
	res := t.DB.Model(&model.Traffic{}).Where(key+" = ?", id).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Traffic: %v", res.Error)
		return 0, false, i18n.CountModelError
	}
	return count, true, i18n.Success
}

func (t *TrafficRepo) GetByKey(key string, id uint, limit, offset int, preloadL ...string) ([]model.Traffic, int64, bool, string) {
	var (
		traffics       = make([]model.Traffic, 0)
		count, ok, msg = t.CountByKey(key, id)
	)
	if !ok {
		return traffics, count, false, msg
	}
	res := t.DB.Model(&model.Traffic{}).Where(key+" = ?", id)
	res = preload(res, preloadL...).Limit(limit).Offset(offset).Find(&traffics)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Traffics: %v", res.Error)
		return traffics, count, false, i18n.GetTrafficError
	}
	return traffics, count, true, i18n.Success
}
