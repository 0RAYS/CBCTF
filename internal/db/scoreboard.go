package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UpdateTeamRanking 更新队伍分数和排名, 更新 redis 排行榜
func UpdateTeamRanking(tx *gorm.DB, contestID uint) (bool, string) {
	var teams []model.Team
	res := tx.Model(&model.Team{}).Where("contest_id = ? AND banned = ?", contestID, false).
		Preload(clause.Associations).Order("score DESC, last ASC").Find(&teams)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get teams: %v", res.Error)
		return false, "GetTeamError"
	}
	for _, team := range teams {
		team.Score, _, _ = CalcTeamScore(tx, team.ContestID, team.ID)
		// 不考虑更新失败的情况
		UpdateTeam(tx, team.ID, map[string]interface{}{"score": team.Score})
	}
	res = tx.Model(&model.Team{}).Where("contest_id = ? AND banned = ?", contestID, false).
		Preload(clause.Associations).Order("score DESC, last ASC").Find(&teams)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get teams: %v", res.Error)
		return false, "GetTeamError"
	}
	err := redis.UpdateTeamRanking(contestID, teams)
	if err != nil {
		log.Logger.Warningf("Failed to update ranking: %v", err)
		return false, "UpdateRankingError"
	}
	return true, "Success"
}

// GetTeamRanking 获取队伍排名
func GetTeamRanking(tx *gorm.DB, contestID uint, limit, offset int) ([]model.Team, int64, bool, string) {
	var count int64
	res := tx.Model(&model.Team{}).Where("contest_id = ? AND banned = ?", contestID, false)
	if err := res.Count(&count).Error; err != nil {
		log.Logger.Warningf("Failed to count teams: %v", err)
		return make([]model.Team, 0), -1, false, "UnknownError"
	}
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	if teams, err := redis.GetTeamRanking(contestID, int64(offset), int64(limit)-1); err == nil && teams != nil {
		return teams, count, true, "Success"
	}
	var teams []model.Team
	res = res.Preload(clause.Associations).Order("score DESC, last ASC").Find(&teams)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get teams: %v", res.Error)
		return make([]model.Team, 0), -1, false, "GetTeamError"
	}
	go UpdateTeamRanking(tx, contestID)
	return teams[offset:limit], count, true, "Success"
}

// UpdateUserRanking 更新用户分数和排名, 更新 redis 排行榜, 主要依赖定时任务调用
func UpdateUserRanking(tx *gorm.DB) (bool, string) {
	var users []model.User
	res := tx.Model(&model.User{}).Where("banned = ?", false).
		Order("score DESC, solved DESC").Find(&users)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get users: %v", res.Error)
		return false, "GetUserError"
	}
	err := redis.UpdateUserRanking(users)
	if err != nil {
		log.Logger.Warningf("Failed to update ranking: %v", err)
		return false, "UpdateRankingError"
	}
	return true, "Success"
}

// GetUserRanking 获取用户排名
func GetUserRanking(tx *gorm.DB, limit, offset int) ([]model.User, int64, bool, string) {
	var count int64
	res := tx.Model(&model.User{}).Where("banned = ?", false)
	if err := res.Count(&count).Error; err != nil {
		log.Logger.Warningf("Failed to count users: %v", err)
		return make([]model.User, 0), -1, false, "UnknownError"
	}
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	if users, err := redis.GetUserRanking(int64(offset), int64(limit)-1); err == nil && users != nil {
		return users, count, true, "Success"
	}
	var users []model.User
	res = tx.Model(&model.User{}).Where("banned = ?", false).Order("score DESC, solved DESC").Find(&users)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get users: %v", res.Error)
		return make([]model.User, 0), -1, false, "UnknownError"
	}
	go UpdateUserRanking(tx)
	return users, count, true, "Success"
}

// GetTeamRankDetail 获取队伍排名详情, 提供额外的分数变动记录
func GetTeamRankDetail(tx *gorm.DB, contestID uint, limit, offset int) ([]map[string]interface{}, bool, string) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	var data []map[string]interface{}
	teams, _, ok, msg := GetTeamRanking(tx, contestID, limit, offset)
	if !ok {
		return data, false, msg
	}
	for _, team := range teams {
		submissions, ok, msg := GetTeamSolved(tx, team.ID)
		if !ok {
			return data, false, msg
		}
		var history []map[string]interface{}
		for _, submission := range submissions {
			history = append(history, map[string]interface{}{
				"challenge": submission.ChallengeID,
				"score":     submission.Score,
				"time":      submission.CreatedAt,
			})
		}
		data = append(data, map[string]interface{}{"team": team, "history": history})
	}
	return data, true, "Success"
}
