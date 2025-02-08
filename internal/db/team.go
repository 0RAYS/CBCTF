package db

import (
	"CBCTF/internal/form"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateTeam 创建队伍, 名称在 model.Contest 中唯一
func CreateTeam(tx *gorm.DB, form form.CreateTeamForm, captain model.User, contest model.Contest) (model.Team, bool, string) {
	if !IsUniqueTeamName(tx, form.Name, contest.ID) {
		return model.Team{}, false, "TeamNameExists"
	}
	if !IsUniqueTeamMember(tx, contest.ID, captain.ID) {
		return model.Team{}, false, "TeamMemberExists"
	}
	team := model.InitTeam(form, captain, contest.ID)
	res := tx.Model(&model.Team{}).Create(&team)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create team: %s", res.Error)
		return model.Team{}, false, "CreateTeamError"
	}
	go func() {
		if err := redis.DelTeamsCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete teams cache: %s", err)
		}
	}()
	return team, true, "Success"
}

// GetTeamByID 根据 ID 获取 model.Team
func GetTeamByID(tx *gorm.DB, id uint, preloadL ...bool) (model.Team, bool, string) {
	preload := true
	nest := false
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if len(preloadL) > 1 {
		nest = preloadL[1]
	}
	if team, ok := redis.GetTeamCache(id, redis.GetType(preload, nest)); ok {
		return team, true, "Success"
	}
	var team model.Team
	res := tx.Model(&model.Team{}).Where("id = ?", id)
	if preload {
		if nest {
			res = res.Preload("Users.Teams").Preload("Users.Contests")
		}
		res = res.Preload(clause.Associations)
	}
	res = res.Find(&team).Limit(1)
	if res.RowsAffected != 1 {
		return model.Team{}, false, "TeamNotFound"
	}
	go func() {
		if err := redis.SetTeamCache(team, redis.GetType(preload, nest)); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to set team cache: %s", err)
		}
	}()
	return team, true, "Success"
}

// GetTeamByName 根据名称获取 model.Team, name 用户可控, 不进行缓存
func GetTeamByName(tx *gorm.DB, name string, contestID uint, preloadL ...bool) (model.Team, bool, string) {
	preload := true
	nest := false
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if len(preloadL) > 1 {
		nest = preloadL[1]
	}
	var team model.Team
	res := tx.Model(&model.Team{}).Where("name = ? AND contest_id = ?", name, contestID)
	if preload {
		if nest {
			res = res.Preload("Users.Teams").Preload("Users.Contests")
		}
		res = res.Preload(clause.Associations)
	}
	res = res.Find(&team).Limit(1)
	if res.RowsAffected != 1 {
		return model.Team{}, false, "TeamNotFound"
	}
	return team, true, "Success"
}

// GetTeamByUserID 根据 UserID 获取 model.Team, 结果等同于 GetTeam preload = true, nest = false
func GetTeamByUserID(tx *gorm.DB, userID uint, contestID uint) (model.Team, bool, string) {
	user, ok, msg := GetUserByID(tx, userID, true, true)
	if !ok {
		return model.Team{}, false, msg
	}
	for _, team := range user.Teams {
		if team.ContestID == contestID {
			return *team, true, "Success"
		}
	}
	return model.Team{}, false, "UserNotInTeam"
}

// DeleteTeam 删除 model.Team, 同时删除与 model.User, model.Contest 的关联
func DeleteTeam(tx *gorm.DB, team model.Team) (bool, string) {
	contest, ok, msg := GetContestByID(tx, team.ContestID)
	if !ok {
		return false, msg
	}
	// 删除 User 和 Contest 关联
	for _, user := range team.Users {
		if err := DeleteUserFromContest(tx, *user, contest); err != nil {
			log.Logger.Warningf("Failed to delete user_contest: %s", err)
			return false, "DeleteUserFromContestError"
		}
	}
	if err := tx.Model(&model.Team{}).Select(clause.Associations).Delete(&team).Error; err != nil {
		log.Logger.Warningf("Failed to delete team: %s", err)

		return false, "DeleteTeamError"
	}
	if !ClearByID(tx, "team_id", team.ID) {
		return false, "DeleteAssociatedDataError"
	}
	go func() {
		if err := redis.DelTeamCache(team.ID); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete team cache: %s", err)
		}
		if err := redis.DelTeamsCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete teams cache: %s", err)
		}
	}()
	return true, "Success"
}

// UpdateTeam 使用 map 更新属性, 使用结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateTeam(tx *gorm.DB, id uint, updateData map[string]interface{}) (bool, string) {
	res := tx.Model(&model.Team{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update team: %s", res.Error)
		return false, "UpdateTeamError"
	}
	go func() {
		if err := redis.DelTeamCache(id); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete team cache: %s", err)
		}
		if err := redis.DelTeamsCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete teams cache: %s", err)
		}
	}()
	return true, "Success"
}

// JoinTeam model.User 加入 model.Team, 建立三个模型直接的关联关系
func JoinTeam(tx *gorm.DB, user model.User, team model.Team, contest model.Contest) (bool, string) {
	if !IsUniqueTeamMember(tx, contest.ID, user.ID) {
		return false, "TeamMemberExists"
	}
	if team.Banned {
		return false, "TeamIsBanned"
	}
	if len(team.Users)+1 > contest.Size {
		return false, "TeamIsFull"
	}
	// 关联 Team User Many2Many
	if err := AppendUserToTeam(tx, user, team); err != nil {
		log.Logger.Warningf("Failed to insert user_team: %s", err)
		return false, "AppendUserToTeamError"
	}
	// 关联 User Contest Many2Many
	if err := AppendUserToContest(tx, user, contest); err != nil {
		log.Logger.Warningf("Failed to insert user_contest: %s", err)
		return false, "AppendContestToUserError"
	}
	go func() {
		if err := redis.DelUserCache(user.ID); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete user cache: %s", err)
		}
		if err := redis.DelUsersCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete users cache: %s", err)
		}
		if err := redis.DelTeamCache(team.ID); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete team cache: %s", err)
		}
		if err := redis.DelTeamsCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete teams cache: %s", err)
		}
		if err := redis.DelContestCache(contest.ID); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete contest cache: %s", err)
		}
		if err := redis.DelContestsCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete contests cache: %s", err)
		}
	}()
	return true, "Success"
}

// LeaveTeam model.User 离开 model.Team, 删除三个模型直接的关联关系
func LeaveTeam(tx *gorm.DB, user model.User, team model.Team, contest model.Contest) (bool, string) {
	if !IsMemberInTeam(tx, team.ID, user.ID) {
		return false, "UserNotInTeam"
	}
	if team.CaptainID == user.ID {
		return false, "CaptainCannotLeave"
	}
	// 队伍人数为 1 时一定是队长, 无法到达这个代码, 暂且保留; 退出后队伍人数为0, 删除队伍;
	if len(team.Users) == 1 {
		DeleteTeam(tx, team)
	}
	if err := DeleteUserFromTeam(tx, user, team); err != nil {
		log.Logger.Warningf("Failed to delete user_team: %s", err)
		return false, "DeleteUserFromTeamError"
	}
	if err := DeleteUserFromContest(tx, user, contest); err != nil {
		log.Logger.Warningf("Failed to delete user_contest: %s", err)
		return false, "DeleteUserFromContestError"
	}
	go func() {
		if err := redis.DelUserCache(user.ID); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete user cache: %s", err)
		}
		if err := redis.DelUsersCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete users cache: %s", err)
		}
		if err := redis.DelTeamCache(team.ID); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete team cache: %s", err)
		}
		if err := redis.DelTeamsCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete teams cache: %s", err)
		}
		if err := redis.DelContestCache(contest.ID); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete contest cache: %s", err)
		}
		if err := redis.DelContestsCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to delete contests cache: %s", err)
		}
	}()
	return true, "Success"
}

// GetTeams 获取 model.Team 列表, preloadL[0] 是否预加载, preloadL[1] 是否嵌套预加载
func GetTeams(tx *gorm.DB, contestID uint, limit int, offset int, all bool, preloadL ...bool) ([]model.Team, int64, bool, string) {
	preload := true
	nest := false
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if len(preloadL) > 1 {
		nest = preloadL[1]
	}
	var teams []model.Team
	var count int64
	res := tx.Model(&model.Team{}).Where("contest_id = ?", contestID)
	if !all {
		res = res.Where("hidden = ? AND banned = ?", false, false)
	}
	if res.Count(&count).Error != nil {
		log.Logger.Warningf("Failed to get contest count: %s", res.Error)
		return nil, 0, false, "UnknownError"
	}
	if teams, ok := redis.GetTeamsCache(redis.GetType(preload, nest)); ok {
		limit, offset = utils.TidyPaginate(len(teams), limit, offset)
		return teams[limit:offset], count, true, "Success"
	}
	if preload {
		if nest {
			res = res.Preload("Users.Teams").Preload("Users.Contests")
		}
		res = res.Preload(clause.Associations)
	}
	if res = res.Find(&teams); res.Error != nil {
		log.Logger.Warningf("Failed to get teams: %s", res.Error)
		return nil, 0, false, "GetTeamError"
	}
	go func() {
		if err := redis.SetTeamsCache(teams, redis.GetType(preload, nest)); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Logger.Warningf("Failed to set teams cache: %s", err)
		}
	}()
	limit, offset = utils.TidyPaginate(len(teams), limit, offset)
	return teams[limit:offset], count, true, "Success"
}
