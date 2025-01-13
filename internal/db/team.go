package db

import (
	"RayWar/internal/log"
	"RayWar/internal/model"
	"gorm.io/gorm"
)

// CreateTeam 创建新队伍，单独判断在指定 model.Contest 中是否唯一
func CreateTeam(contest model.Contest, creator model.User, name string) (model.Team, bool, string) {
	if !isUniqueTeamName(contest.ID, name) {
		return model.Team{}, false, "TeamNameExists"
	}
	if !isNotRepeatPlayer(contest.ID, creator.ID) {
		return model.Team{}, false, "RepeatPlayer"
	}
	team := model.InitTeam(name)
	res := DB.Model(&model.Team{}).Create(&team)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create team: %v", res.Error.Error())
		return model.Team{}, false, "CreateTeamError"
	}
	err := DB.Model(&team).Association("Users").Append(&creator)
	if err != nil {
		log.Logger.Warningf("Failed to create user_team: %v", err.Error())
		return model.Team{}, false, "CreateTeamError"
	}
	err = DB.Model(&team).Association("Contests").Append(&contest)
	if err != nil {
		log.Logger.Warningf("Failed to create team_contest: %v", err.Error())
		return model.Team{}, false, "CreateTeamError"
	}
	return team, true, "Success"
}

// GetTeamByID 对于 model.Team 来说，只有 ID 是唯一的，Name 仅在某一 model.Contest 中唯一
func GetTeamByID(id uint, preloadL ...bool) (model.Team, bool, string) {
	var team model.Team
	var res *gorm.DB
	preload := true
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if preload {
		res = DB.Model(&model.Team{}).Where("id = ?", id).Preload("Users").Preload("Contests").
			Find(&team).Limit(1)
	} else {
		res = DB.Model(&model.Team{}).Where("id = ?", id).Find(&team).Limit(1)
	}
	if res.RowsAffected != 1 {
		return model.Team{}, false, "TeamNotFound"
	}
	return team, true, "Success"
}

// UpdateTeam 对字段值的具体要求应当交给上层实现
func UpdateTeam(team model.Team, updateData map[string]interface{}) (bool, string) {
	res := DB.Model(&model.Team{}).Where("id = ?", team.ID).Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update team: %v", res.Error.Error())
		return false, "UpdateError"
	}
	return true, "Success"
}

// DeleteTeam 同时删除与 model.User model.Contest 的关联关系，但不删除关联的数据
func DeleteTeam(team model.Team) (bool, string) {
	if err := DB.Model(&team).Select("Users", "Contests").Delete(&team).Error; err != nil {
		log.Logger.Warningf("Failed to delete team: %v", err.Error())
		return false, "DeleteTeamError"
	}
	return true, "Success"
}

// ClearEmptyTeam 清除所有人数为0的队伍，目前想法为定时任务执行
func ClearEmptyTeam() []map[string]interface{} {
	var teams []model.Team
	var deleted []map[string]interface{}
	DB.Model(&model.Team{}).Preload("Users").Find(&teams)
	for _, team := range teams {
		if len(team.Users) == 0 {
			if ok, _ := DeleteTeam(team); ok {
				deleted = append(deleted, map[string]interface{}{"name": team.Name, "desc": team.Desc})
			}
		}
	}
	return deleted
}

// GetTeams 获取所有队伍
func GetTeams(limit int, offset int, all bool, contestIDL ...uint) ([]model.Team, int64, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
	var teams []model.Team
	var total int64
	if all {
		if res := DB.Model(&model.Team{}).Count(&total); res.Error != nil {
			log.Logger.Warningf("Failed to get teams: %s", res.Error.Error())
			return nil, 0, false, "UnknownError"
		}
		if res := DB.Model(&model.Team{}).Limit(limit).Offset(offset).Find(&teams); res.Error != nil {
			log.Logger.Warningf("Failed to get teams: %s", res.Error.Error())
			return nil, 0, false, "TeamNotFound"
		}
	} else {
		if res := DB.Model(&model.Team{}).Preload("Contests").Where("id = ?", contestIDL[0]).
			Where("hidden != ? and banned != ?", true, true).Count(&total); res.Error != nil {
			log.Logger.Warningf("Failed to get teams: %s", res.Error.Error())
			return nil, 0, false, "UnknownError"
		}
		if res := DB.Model(&model.Team{}).Preload("Contests").Where("id = ?", contestIDL[0]).
			Where("hidden != ? and banned != ?", true, true).Limit(limit).Offset(offset).Find(&teams); res.Error != nil {
			log.Logger.Warningf("Failed to get teams: %s", res.Error.Error())
			return nil, 0, false, "TeamNotFound"
		}
	}
	return teams, total, true, "Success"
}
