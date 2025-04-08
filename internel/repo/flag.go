package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
	"time"
)

type FlagRepo struct {
	Repo[model.Flag]
}

type CreateFlagOptions struct {
	ContestID    uint
	UsageID      uint
	Value        string
	Score        float64
	CurrentScore float64
	Decay        float64
	MinScore     float64
	ScoreType    uint
	Blood        model.Uints
}

type UpdateFlagOptions struct {
	Value        *string      `json:"value"`
	Score        *float64     `json:"score"`
	CurrentScore *float64     `json:"current_score"`
	Decay        *float64     `json:"decay"`
	MinScore     *float64     `json:"min_score"`
	ScoreType    *uint        `json:"score_type"`
	Solvers      *int64       `json:"solvers"`
	Blood        *model.Uints `json:"blood"`
	Last         *time.Time   `json:"last"`
}

func InitFlagRepo(tx *gorm.DB) *FlagRepo {
	return &FlagRepo{Repo: Repo[model.Flag]{DB: tx, Model: "Flag"}}
}

func (f *FlagRepo) Count(key string, id uint) (int64, bool, string) {
	var count int64
	res := f.DB.Model(&model.Flag{}).Where(key+" = ?", id).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Flags: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (f *FlagRepo) GetByKeyID(key string, id uint, limit, offset int, preloadL ...string) ([]model.Flag, int64, bool, string) {
	var (
		flags          = make([]model.Flag, 0)
		count, ok, msg = f.Count(key, id)
	)
	if !ok {
		return flags, count, false, msg
	}
	res := f.DB.Model(&model.Flag{}).Where(key+" = ?", id)
	res = preload(res, preloadL...).Limit(limit).Offset(offset).Find(&flags)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Flags: %s", res.Error)
		return flags, count, false, "GetFlagError"
	}
	return flags, count, true, "Success"
}

func (f *FlagRepo) Update(id uint, options UpdateFlagOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Flag: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		flag, ok, msg := f.GetByID(id)
		if !ok {
			return false, msg
		}
		data["version"] = flag.Version + 1
		res := f.DB.Model(&model.Flag{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, flag.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Flag: %s", res.Error)
			return false, "UpdateFlagError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}
