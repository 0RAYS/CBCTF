package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateContest 创建比赛
func CreateContest(ctx context.Context, name string) (model.Contest, bool, string) {
	if !isUniqueName(name, model.Contest{}) {
		return model.Contest{}, false, "ContestNameExists"
	}
	contest := model.InitContest(name)
	res := DB.WithContext(ctx).Model(&model.Contest{}).Create(&contest)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create contest: %s", res.Error.Error())
		return model.Contest{}, false, "CreateContestError"
	}
	return contest, true, "Success"
}

// GetContestByID 根据 ID 获取 model.Contest
func GetContestByID(ctx context.Context, id uint, preloadL ...bool) (model.Contest, bool, string) {
	var contest model.Contest
	var res *gorm.DB
	preload := true
	nest := false
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if len(preloadL) > 1 {
		nest = preloadL[1]
	}
	if preload {
		if nest {
			res = DB.WithContext(ctx).Model(&model.Contest{}).Where("id = ?", id).
				Preload("Teams.Users").Preload("Users.Contests").Preload("Users.Teams").
				Preload(clause.Associations).Find(&contest).Limit(1)
		} else {
			res = DB.WithContext(ctx).Model(&model.Contest{}).Where("id = ?", id).Preload(clause.Associations).
				Find(&contest).Limit(1)
		}
	} else {
		res = DB.WithContext(ctx).Model(&model.Contest{}).Where("id = ?", id).Find(&contest).Limit(1)
	}
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
