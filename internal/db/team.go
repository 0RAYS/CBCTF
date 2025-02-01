package db

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"gorm.io/gorm/clause"
)

// CreateTeam 创建队伍, 名称在 model.Contest 中唯一
func CreateTeam(ctx context.Context, form constants.CreateTeamForm, captainID uint, contestID uint) (model.Team, bool, string) {
	if !IsUniqueTeamName(form.Name, contestID) {
		return model.Team{}, false, "TeamNameExists"
	}
	if !IsUniqueTeamMember(contestID, captainID) {
		return model.Team{}, false, "TeamMemberExists"
	}
	team := model.InitTeam(form, captainID)
	res := DB.WithContext(ctx).Model(&model.Team{}).Create(&team)
	if res.Error != nil {
		log.Logger.Warningf("Failed to create team: %s", res.Error)
		return model.Team{}, false, "CreateTeamError"
	}
	if ok, msg := JoinTeam(ctx, captainID, contestID, team.ID); !ok {
		return model.Team{}, false, msg
	}
	//go func() {
	//	if err := redis.DelTeamsCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete teams cache: %s", err)
	//	}
	//}()
	return team, true, "Success"
}

// GetTeamByID 根据 ID 获取 model.Team
func GetTeamByID(ctx context.Context, id uint, preloadL ...bool) (model.Team, bool, string) {
	preload := true
	nest := false
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if len(preloadL) > 1 {
		nest = preloadL[1]
	}
	//cacheKey := fmt.Sprintf("team:%d:%v:%v", id, preload, nest)
	//if team, ok := redis.GetTeamCache(cacheKey); ok {
	//	return team, true, "Success"
	//}
	var team model.Team
	res := DB.WithContext(ctx).Model(&model.Team{}).Where("id = ?", id)
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
	//go func() {
	//	if err := redis.SetTeamCache(cacheKey, team); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to set team cache: %s", err)
	//	}
	//}()
	return team, true, "Success"
}

// GetTeamByName 根据名称获取 model.Team, name 用户可控, 不进行缓存
func GetTeamByName(ctx context.Context, name string, contestID uint, preloadL ...bool) (model.Team, bool, string) {
	preload := true
	nest := false
	if len(preloadL) > 0 {
		preload = preloadL[0]
	}
	if len(preloadL) > 1 {
		nest = preloadL[1]
	}
	var team model.Team
	res := DB.WithContext(ctx).Model(&model.Team{}).Where("name = ? AND contest_id = ?", name, contestID)
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
func GetTeamByUserID(ctx context.Context, userID uint, contestID uint) (model.Team, bool, string) {
	user, ok, msg := GetUserByID(ctx, userID, true, true)
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

// DeleteTeam 根据 id 删除 model.Team, 同时删除与 model.User, model.Contest 的关联
func DeleteTeam(ctx context.Context, id uint) (bool, string) {
	team, ok, msg := GetTeamByID(ctx, id, true)
	if !ok {
		return false, msg
	}
	contest, ok, msg := GetContestByID(ctx, team.ContestID, true)
	if !ok {
		return false, msg
	}
	// 删除 User 和 Contest 关联
	for _, user := range team.Users {
		if err := DeleteUserFromContest(ctx, *user, contest); err != nil {
			log.Logger.Warningf("Failed to delete user_contest: %s", err)
			return false, "DeleteUserFromContestError"
		}
	}
	if err := DB.WithContext(ctx).Model(&model.Team{}).Select(clause.Associations).Delete(&team).Error; err != nil {
		log.Logger.Warningf("Failed to delete team: %s", err)
		return false, "DeleteTeamError"
	}
	//go func() {
	//	if err := redis.DelTeamCache(id); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete team cache: %s", err)
	//	}
	//	if err := redis.DelTeamsCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete teams cache: %s", err)
	//	}
	//}()
	return true, "Success"
}

// UpdateTeam 使用 map 更新属性, 使用结构体会导致零值未更新, 对字段值的具体要求应当交给上层实现
func UpdateTeam(ctx context.Context, id uint, updateData map[string]interface{}) (bool, string) {
	res := DB.WithContext(ctx).Model(&model.Team{}).Where("id = ?", id).
		Omit("id", "created_at", "updated_at", "deleted_at").Updates(updateData)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update team: %s", res.Error)
		return false, "UpdateTeamError"
	}
	//go func() {
	//	if err := redis.DelTeamCache(id); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete team cache: %s", err)
	//	}
	//	if err := redis.DelTeamsCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete teams cache: %s", err)
	//	}
	//}()
	return true, "Success"
}

// JoinTeam model.User 加入 model.Team, 建立三个模型直接的关联关系
func JoinTeam(ctx context.Context, userID uint, contestID uint, teamID uint) (bool, string) {
	user, ok, msg := GetUserByID(ctx, userID, false)
	if !ok {
		return false, msg
	}
	if !IsUniqueTeamMember(contestID, userID) {
		return false, "TeamMemberExists"
	}
	contest, ok, msg := GetContestByID(ctx, contestID, false)
	if !ok || contest.Hidden {
		return false, msg
	}
	team, ok, msg := GetTeamByID(ctx, teamID, true)
	if !ok {
		return false, msg
	}
	if team.Banned {
		return false, "TeamIsBanned"
	}
	if len(team.Users)+1 > contest.Size {
		return false, "TeamIsFull"
	}
	// 关联 Team User Many2Many
	if err := AppendUserToTeam(ctx, user, team); err != nil {
		log.Logger.Warningf("Failed to insert user_team: %s", err)
		return false, "AppendUserToTeamError"
	}
	// 关联 Contest Team HasMany
	if err := AppendTeamToContest(ctx, team, contest); err != nil {
		log.Logger.Warningf("Failed to insert contest_team: %s", err)
		return false, "AppendTeamToContestError"
	}
	// 关联 User Contest Many2Many
	if err := AppendUserToContest(ctx, user, contest); err != nil {
		log.Logger.Warningf("Failed to insert user_contest: %s", err)
		return false, "AppendContestToUserError"
	}
	//go func() {
	//	if err := redis.DelUserCache(userID); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete user cache: %s", err)
	//	}
	//	if err := redis.DelUsersCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete users cache: %s", err)
	//	}
	//	if err := redis.DelTeamCache(teamID); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete team cache: %s", err)
	//	}
	//	if err := redis.DelTeamsCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete teams cache: %s", err)
	//	}
	//	if err := redis.DelContestCache(contestID); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete contest cache: %s", err)
	//	}
	//	if err := redis.DelContestsCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete contests cache: %s", err)
	//	}
	//}()
	return ok, "Success"
}

// LeaveTeam model.User 离开 model.Team, 删除三个模型直接的关联关系
func LeaveTeam(ctx context.Context, userID uint, contestID uint, teamID uint) (bool, string) {
	user, ok, msg := GetUserByID(ctx, userID, false)
	if !ok {
		return false, msg
	}
	team, ok, msg := GetTeamByID(ctx, teamID, true)
	if !ok {
		return false, msg
	}
	if !IsMemberInTeam(team.ID, user.ID) {
		return false, "UserNotInTeam"
	}
	contest, ok, msg := GetContestByID(ctx, contestID, false)
	if !ok {
		return false, msg
	}
	if len(team.Users) > 1 && team.CaptainID == userID {
		return false, "CaptainCannotLeave"
	}
	// 退出后队伍人数为0, 删除队伍
	if len(team.Users) == 1 {
		DeleteTeam(ctx, team.ID)
	}
	if err := DeleteUserFromTeam(ctx, user, team); err != nil {
		log.Logger.Warningf("Failed to delete user_team: %s", err)
		return false, "DeleteUserFromTeamError"
	}
	if err := DeleteUserFromContest(ctx, user, contest); err != nil {
		log.Logger.Warningf("Failed to delete user_contest: %s", err)
		return false, "DeleteUserFromContestError"
	}
	//go func() {
	//	if err := redis.DelUserCache(userID); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete user cache: %s", err)
	//	}
	//	if err := redis.DelUsersCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete users cache: %s", err)
	//	}
	//	if err := redis.DelTeamCache(teamID); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete team cache: %s", err)
	//	}
	//	if err := redis.DelTeamsCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete teams cache: %s", err)
	//	}
	//	if err := redis.DelContestCache(contestID); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete contest cache: %s", err)
	//	}
	//	if err := redis.DelContestsCache(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to delete contests cache: %s", err)
	//	}
	//}()
	return true, "Success"
}

func GetTeams(ctx context.Context, contestID uint, limit int, offset int, all bool, preloadL ...bool) ([]model.Team, int64, bool, string) {
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
	var teams []model.Team
	var count int64
	res := DB.WithContext(ctx).Model(&model.Team{}).Where("contest_id = ?", contestID)
	if !all {
		res = res.Where("hidden = ? AND banned = ?", false, false)
	}
	if res.Count(&count).Error != nil {
		log.Logger.Warningf("Failed to get contest count: %s", res.Error)
		return nil, 0, false, "UnknownError"
	}
	//cacheKey := fmt.Sprintf("teams:%d:%v:%v:%d:%d", contestID, all, preload, limit, offset)
	//if teams, ok := redis.GetTeamsCache(cacheKey); ok {
	//	return teams, count, true, "Success"
	//}
	if preload {
		if nest {
			res = res.Preload("Users.Teams").Preload("Users.Contests")
		}
		res = res.Preload(clause.Associations)
	}
	if res = res.Limit(limit).Offset(offset).Find(&teams); res.Error != nil {
		log.Logger.Warningf("Failed to get teams: %s", res.Error)
		return nil, 0, false, "UnknownError"
	}
	//go func() {
	//	if err := redis.SetTeamsCache(cacheKey, teams); err != nil && !errors.Is(err, context.DeadlineExceeded) {
	//		log.Logger.Warningf("Failed to set teams cache: %s", err)
	//	}
	//}()
	return teams, count, true, "Success"
}
