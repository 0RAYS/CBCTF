package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"gorm.io/gorm/clause"
	"time"
)

// CreateContest 创建比赛
func CreateContest(ctx context.Context, name string, desc string, size int, start time.Time, duration time.Duration, hidden bool) (model.Contest, bool, string) {
	if !IsUniqueName(name, model.Contest{}) {
		return model.Contest{}, false, "ContestNameExists"
	}
	contest := model.InitContest(name, desc, size, start, duration, hidden)
	res := DB.WithContext(ctx).Model(&model.Contest{}).Create(&contest)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create contest: %s", res.Error.Error())
		return model.Contest{}, false, "CreateContestError"
	}
	return contest, true, "Success"
}

// GetContestByID 根据 ID 获取 model.Contest
func GetContestByID(ctx context.Context, id uint, preloadL ...bool) (model.Contest, bool, string) {
	preload := true
	nest := false
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if len(preloadL) > 1 {
		nest = preloadL[1]
	}
	var contest model.Contest
	res := DB.WithContext(ctx).Model(&model.Contest{}).Where("id = ?", id)
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

// DeleteContest 根据 id 删除 model.Contest, 同时删除与 model.Team, model.User 的关联, 同时删除 model.Team
func DeleteContest(ctx context.Context, id uint) (bool, string) {
	contest, ok, msg := GetContestByID(ctx, id)
	if !ok {
		return false, msg
	}
	for _, team := range contest.Teams {
		if ok, msg := DeleteTeam(ctx, team.ID); !ok {
			return false, msg
		}
	}
	if err := DB.WithContext(ctx).Model(&model.Contest{}).Select(clause.Associations).Delete(&contest).Error; err != nil {
		log.Logger.Warningf("Failed to delete contest: %s", err.Error())
		return false, "DeleteContestError"
	}
	return true, "Success"
}

// UpdateContest 使用 map 更新属性, 结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateContest(ctx context.Context, id uint, updateData map[string]interface{}) (bool, string) {
	res := DB.WithContext(ctx).Model(&model.Contest{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update contest: %v", res.Error.Error())
		return false, "UpdateError"
	}
	return true, "Success"
}

func CountContests(ctx context.Context) int64 {
	var count int64
	DB.WithContext(ctx).Model(&model.Contest{}).Count(&count)
	return count
}

func GetContests(ctx context.Context, limit int, offset int, all bool, preloadL ...bool) ([]model.Contest, int64, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
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
	res := DB.WithContext(ctx).Model(&model.Contest{})
	if !all {
		res = res.Where("hidden = ?", false)
	}
	if res.Count(&count).Error != nil {
		log.Logger.Errorf("Failed to get contest count: %s", res.Error.Error())
		return nil, 0, false, "UnknownError"
	}
	if preload {
		if nest {
			res = res.Preload("Teams.Users").Preload("Users.Contests").Preload("Users.Teams")
		}
		res = res.Preload(clause.Associations)
	}
	if res = res.Order("Start desc").Limit(limit).Offset(offset).Find(&contests); res.Error != nil {
		log.Logger.Errorf("Failed to get contests: %s", res.Error.Error())
		return nil, 0, false, "UnknownError"
	}
	return contests, count, true, "Success"

}
