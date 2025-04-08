package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

type TrafficRepo struct {
	Repo[model.Traffic]
}

type CreateTrafficOptions struct {
	SrcIP       string
	DstIP       string
	SrcPort     uint16
	DstPort     uint16
	Payload     string
	Type        string
	ContainerID uint
}

func InitTrafficRepo(tx *gorm.DB) *TrafficRepo {
	return &TrafficRepo{Repo: Repo[model.Traffic]{DB: tx, Model: "Traffic"}}
}

func (t *TrafficRepo) Count(containerID uint) (int64, bool, string) {
	var count int64
	res := t.DB.Model(&model.Traffic{}).Where("container_id = ?", containerID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Traffic: %v", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (t *TrafficRepo) GetAll(containerID uint, limit, offset int, preloadL ...string) ([]model.Traffic, int64, bool, string) {
	var (
		traffics       = make([]model.Traffic, 0)
		count, ok, msg = t.Count(containerID)
	)
	if !ok {
		return traffics, count, false, msg
	}
	res := t.DB.Model(&model.Traffic{}).Where("container_id = ?", containerID)
	res = preload(res, preloadL...).Limit(limit).Offset(offset).Find(&traffics)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Traffics: %v", res.Error)
		return traffics, count, false, "GetTrafficError"
	}
	return traffics, count, true, "Success"
}
