package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

type UsageRepo struct {
	Repo[model.Usage]
}

type CreateUsageOptions struct {
	ContestID   uint
	ChallengeID string
	Name        string
	Desc        string
}

type UpdateUsageOptions struct {
	Name *string
	Desc *string
}

func InitUsageRepo(tx *gorm.DB) *UsageRepo {
	return &UsageRepo{Repo: Repo[model.Usage]{DB: tx, Model: "Usage"}}
}

func (u *UsageRepo) IsUniqueChallenge(contestID uint, challengeID string) bool {
	res := u.DB.Model(&model.Usage{}).Where("contest_id = ? AND challenge_id = ?", contestID, challengeID).
		Find(&model.Usage{}).Limit(1)
	return res.RowsAffected == 0
}

//func (u *UsageRepo) Create(options CreateUsageOptions) (model.Usage, bool, string) {
//	usage, err := utils.S2S[model.Usage](options)
//	if err != nil {
//		log.Logger.Warningf("Failed to convert options to model.Usage: %s", err)
//		return model.Usage{}, false, "Options2ModelError"
//	}
//	if res := u.DB.Model(&model.Usage{}).Create(&usage); res.Error != nil {
//		log.Logger.Warningf("Failed to create Usage: %s", res.Error)
//		return model.Usage{}, false, "CreateUsageError"
//	}
//	return usage, true, "Success"
//}

//func (u *UsageRepo) getByUniqueKey(key string, value interface{}, preload bool, depth int) (model.Usage, bool, string) {
//	switch key {
//	case "id":
//		value = value.(uint)
//	default:
//		return model.Usage{}, false, "UnsupportedKey"
//	}
//	var usage model.Usage
//	res := u.DB.Model(&model.Usage{}).Where(key+" = ?", value)
//	res = model.GetPreload(res, model.Submission{}, preload, depth).Find(&usage).Limit(1)
//	if res.RowsAffected == 0 {
//		return model.Usage{}, false, "UsageNotFound"
//	}
//	return usage, true, "Success"
//}

//func (u *UsageRepo) GetByID(id uint, preload bool, depth int) (model.Usage, bool, string) {
//	return u.getByUniqueKey("id", id, preload, depth)
//}

func (u *UsageRepo) GetBy2ID(contestID uint, challengeID string, preload bool, depth int, hidden bool) (model.Usage, bool, string) {
	var usage model.Usage
	res := u.DB.Model(&model.Usage{}).Where("contest_id = ? AND challenge_id = ?", contestID, challengeID)
	if !hidden {
		res = res.Where("hidden = ?", false)
	}
	res = model.GetPreload(res, u.Model, preload, depth).Find(&usage).Limit(1)
	if res.RowsAffected == 0 {
		return model.Usage{}, false, "UsageNotFound"
	}
	return usage, true, "Success"
}

func (u *UsageRepo) Count(contestID uint, hidden bool) (int64, bool, string) {
	var count int64
	res := u.DB.Model(&model.Usage{}).Where("contest_id = ?", contestID)
	if !hidden {
		res = res.Where("hidden = ?", false)
	}
	res = res.Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Usages: %s", res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (u *UsageRepo) GetAll(contestID uint, limit, offset int, preload bool, depth int, hidden bool) ([]model.Usage, int64, bool, string) {
	var (
		usages         = make([]model.Usage, 0)
		count, ok, msg = u.Count(contestID, hidden)
	)
	if !ok {
		return usages, count, false, msg
	}
	res := u.DB.Model(&model.Usage{}).Where("contest_id = ?", contestID)
	if !hidden {
		res = res.Where("hidden = ?", false)
	}
	res = model.GetPreload(res, u.Model, preload, depth).Find(&usages).Limit(limit).Offset(offset)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get Usages: %s", res.Error)
		return usages, count, false, "GetUsageError"
	}
	return usages, count, true, "Success"
}

func (u *UsageRepo) Update(id uint, options UpdateUsageOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Usage: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		usage, ok, msg := u.GetByID(id, false, 0)
		if !ok {
			return false, msg
		}
		data["version"] = usage.Version + 1
		res := u.DB.Model(&model.Usage{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, usage.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Usage: %s", res.Error)
			return false, "UpdateUsageError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}

//func (u *UsageRepo) Delete(idL ...uint) (bool, string) {
//	res := u.DB.Model(&model.Usage{}).Where("id IN ?", idL).Delete(&model.Usage{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Usage: %s", res.Error)
//		return false, "DeleteUsageError"
//	}
//	return true, "Success"
//}
