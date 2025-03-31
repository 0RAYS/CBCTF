package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
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
	Attempt      int64
}

type UpdateFlagOptions struct {
	Value        *string
	Score        *float64
	CurrentScore *float64
	Decay        *float64
	MinScore     *float64
	ScoreType    *uint
	Attempt      *int64
	Solvers      *int64
}

func InitFlagRepo(tx *gorm.DB) *FlagRepo {
	return &FlagRepo{Repo: Repo[model.Flag]{DB: tx, Model: "Flag"}}
}

//func (f *FlagRepo) Create(options CreateFlagOptions) (model.Flag, bool, string) {
//	flag, err := utils.S2S[model.Flag](options)
//	if err != nil {
//		log.Logger.Warningf("Failed to convert options to model.Flag: %s", err)
//		return model.Flag{}, false, "Options2ModelError"
//	}
//	res := f.DB.Model(&model.Flag{}).Create(&flag)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to create Flag: %s", res.Error)
//		return model.Flag{}, false, "CreateFlagError"
//	}
//	return flag, true, "Success"
//}

//func (f *FlagRepo) getByUniqueKey(key string, value interface{}, preload bool, depth int) (model.Flag, bool, string) {
//	switch key {
//	case "id":
//		value = value.(uint)
//	default:
//		return model.Flag{}, false, "UnsupportedKey"
//	}
//	var flag model.Flag
//	res := f.DB.Model(&model.Flag{}).Where(key+" = ?", value).First(&flag)
//	res = model.GetPreload(res, model.Flag{}, preload, depth).Find(&flag).Limit(1)
//	if res.RowsAffected == 0 {
//		return model.Flag{}, false, "FlagNotFound"
//	}
//	return flag, true, "Success"
//}

//func (f *FlagRepo) GetByID(id uint, preload bool, depth int) (model.Flag, bool, string) {
//	return f.getByUniqueKey("id", id, preload, depth)
//}

func (f *FlagRepo) Count(key string, id uint) (int64, bool, string) {
	var count int64
	res := f.DB.Model(&model.Flag{}).Where(key+" = ?", id).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Flags: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (f *FlagRepo) GetByKeyID(key string, id uint, limit, offset int, preload bool, depth int) ([]model.Flag, int64, bool, string) {
	var (
		flags          = make([]model.Flag, 0)
		count, ok, msg = f.Count(key, id)
	)
	if !ok {
		return flags, count, false, msg
	}
	res := f.DB.Model(&model.Flag{}).Where(key+" = ?", id)
	res = model.GetPreload(res, f.Model, preload, depth).Find(&flags).Limit(limit).Offset(offset)
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
		flag, ok, msg := f.GetByID(id, false, 0)
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

//func (f *FlagRepo) Delete(idL ...uint) (bool, string) {
//	res := f.DB.Model(&model.Flag{}).Where("id IN ?", idL).Delete(&model.Flag{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Flag: %s", res.Error)
//		return false, "DeleteFlagError"
//	}
//	return true, "Success"
//}
