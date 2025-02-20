package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpdateRanking(tx *gorm.DB, contestID uint) (bool, string) {
	if !config.Env.Redis.On {
		return false, "RedisOff"
	}
	var teams []model.Team
	res := tx.Model(&model.Team{}).Where("contest_id = ? AND banned = ?", contestID, false).
		Preload(clause.Associations).Order("score DESC, last ASC").Find(&teams)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get teams: %v", res.Error)
		return false, "GetTeamError"
	}
	err := redis.UpdateRanking(contestID, teams)
	if err != nil {
		log.Logger.Warningf("Failed to update ranking: %v", err)
		return false, "UpdateRankingError"
	}
	return true, "Success"
}

func GetRanking(contestID uint, limit, offset int) ([]model.Team, int64, bool, string) {
	var count int64
	res := DB.Model(&model.Team{}).Where("contest_id = ? AND banned = ?", contestID, false)
	if err := res.Count(&count).Error; err != nil {
		log.Logger.Warningf("Failed to count teams: %v", err)
		return nil, -1, false, "UnknownError"
	}
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	if teams, err := redis.GetCachedRanking(contestID, int64(limit), int64(offset)); err == nil && teams != nil {
		return teams, count, true, "Success"
	}
	var teams []model.Team
	res = res.Preload(clause.Associations).Order("score DESC, last ASC").Find(&teams)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get teams: %v", res.Error)
		return nil, -1, false, "GetTeamError"
	}
	go UpdateRanking(DB, contestID)
	return teams[offset:limit], count, true, "Success"
}

func GetRankDetail(contestID uint) ([]map[string]interface{}, bool, string) {
	var data []map[string]interface{}
	teams, _, ok, msg := GetRanking(contestID, 10, 0)
	if !ok {
		return data, false, msg
	}
	for _, team := range teams {
		submissions, ok, msg := GetTeamSolved(DB, contestID, team.ID)
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
