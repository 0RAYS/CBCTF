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
	for _, team := range teams {
		team.Score, _, _ = CalcTeamScore(tx, team.ContestID, team.ID)
		go UpdateTeam(tx, team.ID, map[string]interface{}{"score": team.Score})
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
		return make([]model.Team, 0), -1, false, "UnknownError"
	}
	limit, offset = utils.TidyPaginate(int(count), limit, offset)
	if teams, err := redis.GetCachedRanking(contestID, int64(offset), int64(limit)-1); err == nil && teams != nil {
		return teams, count, true, "Success"
	}
	var teams []model.Team
	res = res.Preload(clause.Associations).Order("score DESC, last ASC").Find(&teams)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get teams: %v", res.Error)
		return make([]model.Team, 0), -1, false, "GetTeamError"
	}
	go UpdateRanking(DB, contestID)
	return teams[offset:limit], count, true, "Success"
}

func GetRankDetail(contestID uint, limit, offset int) ([]map[string]interface{}, bool, string) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	var data []map[string]interface{}
	teams, _, ok, msg := GetRanking(contestID, limit, offset)
	if !ok {
		return data, false, msg
	}
	for _, team := range teams {
		submissions, ok, msg := GetTeamSolved(DB, team.ID)
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
