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
	Attempt     int64
	Hints       model.Strings
	Tags        model.Strings
	Docker      model.Docker
	Dockers     model.Dockers
}

type UpdateUsageOptions struct {
	Name    *string        `json:"name"`
	Desc    *string        `json:"desc"`
	Attempt *int64         `json:"attempt"`
	Hidden  *bool          `json:"hidden"`
	Hints   *model.Strings `json:"hints"`
	Tags    *model.Strings `json:"tags"`
	Docker  *model.Docker  `json:"docker"`
	Dockers *model.Dockers `json:"dockers"`
}

func InitUsageRepo(tx *gorm.DB) *UsageRepo {
	return &UsageRepo{Repo: Repo[model.Usage]{DB: tx, Model: "Usage"}}
}

func (u *UsageRepo) IsUniqueChallenge(contestID uint, challengeID string) bool {
	res := u.DB.Model(&model.Usage{}).Where("contest_id = ? AND challenge_id = ?", contestID, challengeID).
		Find(&model.Usage{}).Limit(1).Find(&model.Usage{})
	return res.RowsAffected == 0
}

func (u *UsageRepo) GetBy2ID(contestID uint, challengeID string, hidden bool, preloadL ...string) (model.Usage, bool, string) {
	var usage model.Usage
	res := u.DB.Model(&model.Usage{}).Where("contest_id = ? AND challenge_id = ?", contestID, challengeID)
	if !hidden {
		res = res.Where("hidden = ?", false)
	}
	res = preload(res, preloadL...).Limit(1).Find(&usage)
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

func (u *UsageRepo) GetAll(contestID uint, limit, offset int, hidden bool, preloadL ...string) ([]model.Usage, int64, bool, string) {
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
	res = preload(res, preloadL...).Limit(limit).Offset(offset).Find(&usages)
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
		usage, ok, msg := u.GetByID(id)
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
