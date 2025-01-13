package db

import (
	"RayWar/internal/log"
	"RayWar/internal/model"
	"gorm.io/gorm"
)

// CreateContest 创建新赛事，判断赛事名是否唯一
func CreateContest(name string) (model.Contest, bool, string) {
	if !isUniqueContestName(name) {
		return model.Contest{}, false, "ContestNameExists"
	}
	team := model.InitContest(name)
	res := DB.Model(&model.Contest{}).Create(&team)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create contest: %s", res.Error.Error())
		return model.Contest{}, false, "CreateContestError"
	}
	return team, true, "Success"
}

// GetContestByName 根据 Name 获取 model.Contest
func GetContestByName(name string, preloadL ...bool) (model.Contest, bool, string) {
	var contest model.Contest
	var res *gorm.DB
	preload := true
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if preload {
		res = DB.Model(&model.Contest{}).Where("name = ?", name).Preload("Teams").
			Find(&contest).Limit(1)
	} else {
		res = DB.Model(&model.Contest{}).Where("name = ?", name).Find(&contest).Limit(1)
	}
	if res.RowsAffected != 1 {
		return model.Contest{}, false, "ContestNotFound"
	}
	return contest, true, "Success"
}

// GetContestByID 根据 ID 获取 model.Contest
func GetContestByID(id uint, preloadL ...bool) (model.Contest, bool, string) {
	var contest model.Contest
	var res *gorm.DB
	preload := true
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if preload {
		res = DB.Model(&model.Contest{}).Where("id = ?", id).Preload("Teams").
			Find(&contest).Limit(1)
	} else {
		res = DB.Model(&model.Contest{}).Where("id = ?", id).Find(&contest).Limit(1)
	}
	if res.RowsAffected != 1 {
		return model.Contest{}, false, "ContestNotFound"
	}
	return contest, true, "Success"
}

// UpdateContest 对字段的具体要求应当交给上层实现
func UpdateContest(contest model.Contest, updateData map[string]interface{}) (bool, string) {
	res := DB.Model(&model.Contest{}).Where("id = ?", contest.ID).Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update contest: %s", res.Error.Error())
		return false, "UpdateError"
	}
	return true, "Success"
}

// DeleteContest 删除该 model.Contest 下所有 model.Team，同时删除与 model.Team 的关联关系
func DeleteContest(contest model.Contest) (bool, string) {
	for _, team := range contest.Teams {
		DeleteTeam(*team)
	}
	if err := DB.Model(&model.Contest{}).Select("Teams").Delete(&contest).Error; err != nil {
		log.Logger.Warningf("Failed to delete contest: %s", err.Error())
		return false, "DeleteContestError"
	}
	return true, "Success"
}

// GetContests 获取所有赛事
func GetContests(limit int, offset int, all bool) ([]model.Contest, int64, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
	var contests []model.Contest
	var total int64
	if all {
		if res := DB.Model(&model.Contest{}).Count(&total); res.Error != nil {
			log.Logger.Warningf("Failed to get contests: %s", res.Error.Error())
			return nil, 0, false, "UnknownError"
		}
		if res := DB.Model(&model.Contest{}).Limit(limit).Offset(offset).Find(&contests); res.Error != nil {
			log.Logger.Warningf("Failed to get contests: %s", res.Error.Error())
			return nil, 0, false, "ContestNotFound"
		}
	} else {
		if res := DB.Model(&model.Contest{}).Where("hidden != ?", true).
			Count(&total); res.Error != nil {
			log.Logger.Warningf("Failed to get contests: %s", res.Error.Error())
			return nil, 0, false, "UnknownError"
		}
		if res := DB.Model(&model.Contest{}).Where("hidden != ?", true).
			Limit(limit).Offset(offset).Find(&contests); res.Error != nil {
			log.Logger.Warningf("Failed to get contests: %s", res.Error.Error())
			return nil, 0, false, "ContestNotFound"
		}
	}
	return contests, total, true, "Success"
}
