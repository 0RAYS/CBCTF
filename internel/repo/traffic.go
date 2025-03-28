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

//func (t *TrafficRepo) Create(options CreateTrafficOptions) (model.Traffic, bool, string) {
//	traffic, err := utils.S2S[model.Traffic](options)
//	if err != nil {
//		log.Logger.Warningf("Failed to convert options to model.Traffic: %v", err)
//		return model.Traffic{}, false, "Options2ModelError"
//	}
//	res := t.DB.Model(&model.Traffic{}).Create(&traffic)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to create Traffic: %v", res.Error)
//		return model.Traffic{}, false, "CreateTrafficError"
//	}
//	return traffic, true, "Success"
//}

//func (t *TrafficRepo) getByUniqueKey(key string, value interface{}, preload bool, depth int) (model.Traffic, bool, string) {
//	switch key {
//	case "id":
//		value = value.(uint)
//	default:
//		return model.Traffic{}, false, "UnsupportedKey"
//	}
//	var traffic model.Traffic
//	res := t.DB.Model(&model.Traffic{}).Where(key+" = ?", value)
//	res = model.GetPreload(res, model.Traffic{}, preload, depth).Find(&traffic).Limit(1)
//	if res.RowsAffected == 0 {
//		return model.Traffic{}, false, "TrafficNotFound"
//	}
//	return traffic, true, "Success"
//}

//func (t *TrafficRepo) GetByID(id uint, preload bool, depth int) (model.Traffic, bool, string) {
//	return t.getByUniqueKey("id", id, preload, depth)
//}

func (t *TrafficRepo) Count(containerID uint) (int64, bool, string) {
	var count int64
	res := t.DB.Model(&model.Traffic{}).Where("container_id = ?", containerID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Traffic: %v", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (t *TrafficRepo) GetAll(containerID uint, limit, offset int, preload bool, depth int) ([]model.Traffic, int64, bool, string) {
	var (
		traffics       = make([]model.Traffic, 0)
		count, ok, msg = t.Count(containerID)
	)
	if !ok {
		return traffics, count, false, msg
	}
	res := t.DB.Model(&model.Traffic{}).Where("container_id = ?", containerID)
	res = model.GetPreload(res, t.Model, preload, depth).Find(&traffics).Limit(limit).Offset(offset)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Traffics: %v", res.Error)
		return traffics, count, false, "GetTrafficError"
	}
	return traffics, count, true, "Success"
}

//func (t *TrafficRepo) Delete(idL ...uint) (bool, string) {
//	res := t.DB.Model(&model.Traffic{}).Where("id IN ?", idL).Delete(&model.Traffic{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Traffic: %v", res.Error)
//		return false, "DeleteTrafficError"
//	}
//	return true, "Success"
//}
