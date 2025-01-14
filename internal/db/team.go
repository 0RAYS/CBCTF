package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateTeam 创建队伍，名称在 model.Contest 中唯一
func CreateTeam(ctx context.Context, name string, captain model.User, contest model.Contest) (model.Team, bool, string) {
	if !isUniqueTeamName(name, contest) {
		return model.Team{}, false, "TeamNameExists"
	}
	team := model.InitTeam(name, captain, contest)
	res := DB.WithContext(ctx).Model(&model.Team{}).Create(&team)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create team: %s", res.Error.Error())
		return model.Team{}, false, "CreateTeamError"
	}
	return team, true, "Success"
}

// GetTeamByID 根据 ID 获取 model.Team
func GetTeamByID(ctx context.Context, id uint, preloadL ...bool) (model.Team, bool, string) {
	var team model.Team
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
			res = DB.WithContext(ctx).Model(&model.Team{}).Where("id = ?", id).
				Preload("Users.Teams").Preload("Users.Contests").Preload(clause.Associations).
				Find(&team).Limit(1)
		} else {
			res = DB.WithContext(ctx).Model(&model.Team{}).Where("id = ?", id).Preload(clause.Associations).
				Find(&team).Limit(1)
		}
	} else {
		res = DB.WithContext(ctx).Model(&model.Team{}).Where("id = ?", id).Find(&team).Limit(1)
	}
	if res.RowsAffected != 1 {
		return model.Team{}, false, "TeamNotFound"
	}
	return team, true, "Success"
}

// DeleteTeam 根据 id 删除 model.Team，同时删除与 model.User, model.Contest 的关联
func DeleteTeam(ctx context.Context, id uint) (bool, string) {
	if err := DB.WithContext(ctx).Model(&model.Team{}).Select(clause.Associations).Delete(&model.Team{}, id).Error; err != nil {
		log.Logger.Warningf("Failed to delete team: %s", err.Error())
		return false, "DeleteTeamError"
	}
	return true, "Success"
}
