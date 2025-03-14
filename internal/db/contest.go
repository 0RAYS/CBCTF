package db

import (
	"CBCTF/internal/form"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateContest 创建比赛
func CreateContest(tx *gorm.DB, form form.CreateContestForm) (model.Contest, bool, string) {
	if !IsUniqueName(tx, form.Name, model.Contest{}) {
		return model.Contest{}, false, "ContestNameExists"
	}
	contest := model.InitContest(form)
	res := tx.Model(&model.Contest{}).Create(&contest)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create contest: %s", res.Error)
		return model.Contest{}, false, "CreateContestError"
	}
	return contest, true, "Success"
}

// GetContestByID 根据 ID 获取 model.Contest
func GetContestByID(tx *gorm.DB, id uint, preloadL ...bool) (model.Contest, bool, string) {
	preload := true
	nest := false
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if len(preloadL) > 1 {
		nest = preloadL[1]
	}
	var contest model.Contest
	res := tx.Model(&model.Contest{}).Where("id = ?", id)
	if preload {
		if nest {
			res = res.Preload("Teams.Users").Preload("Users.Contests").Preload("Users.Teams")
		}
		res = res.Preload(clause.Associations)
	}
	res = res.Find(&contest).Limit(1)
	if res.RowsAffected != 1 {
		return model.Contest{}, false, "ContestNotFound"
	}
	return contest, true, "Success"
}

// DeleteContest 删除 model.Contest, 同时删除与 model.Team, model.User 的关联, 同时删除 model.Team
func DeleteContest(tx *gorm.DB, contest model.Contest) (bool, string) {
	for _, team := range contest.Teams {
		if ok, msg := DeleteTeam(tx, *team); !ok {
			return false, msg
		}
	}
	if err := tx.Model(&model.Contest{}).Select(clause.Associations).Delete(&contest).Error; err != nil {
		log.Logger.Warningf("Failed to delete contest: %s", err)
		return false, "DeleteContestError"
	}
	if !ClearByID(tx, "contest_id", contest.ID) {
		return false, "DeleteAssociatedDataError"
	}
	return true, "Success"
}

// UpdateContest 使用 map 更新属性, 结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateContest(tx *gorm.DB, id uint, updateData map[string]interface{}) (bool, string) {
	for {
		var contest model.Contest
		res := tx.Model(&model.Contest{}).Where("id = ?", id).Find(&contest).Limit(1)
		if res.RowsAffected != 1 {
			return false, "ContestNotFound"
		}
		res = tx.Model(&contest).Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update contest: %v", res.Error)
			return false, "UpdateContestError"
		}
		if res.RowsAffected == 0 {
			log.Logger.Debug("Failed to update contest due to optimistic lock")
			continue
		}
		break
	}
	return true, "Success"
}

// CountContests 获取比赛数量
func CountContests(tx *gorm.DB) int64 {
	var count int64
	tx.Model(&model.Contest{}).Count(&count)
	return count
}

// GetContests 获取比赛列表
func GetContests(tx *gorm.DB, limit int, offset int, all bool, preloadL ...bool) ([]model.Contest, int64, bool, string) {
	preload := true
	nest := false
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if len(preloadL) > 1 {
		nest = preloadL[1]
	}
	var contests []model.Contest
	var count int64
	res := tx.Model(&model.Contest{})
	if !all {
		res = res.Where("hidden = ?", false)
	}
	if res.Count(&count).Error != nil {
		log.Logger.Warningf("Failed to get contest count: %s", res.Error)
		return make([]model.Contest, 0), 0, false, "UnknownError"
	}
	if preload {
		if nest {
			res = res.Preload("Teams.Users").Preload("Users.Contests").Preload("Users.Teams")
		}
		res = res.Preload(clause.Associations)
	}
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	if res = res.Order("Start DESC").Limit(limit).Offset(offset).Find(&contests); res.Error != nil {
		log.Logger.Warningf("Failed to get contests: %s", res.Error)
		return make([]model.Contest, 0), 0, false, "UnknownError"
	}
	return contests, count, true, "Success"
}
