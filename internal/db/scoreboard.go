package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"gorm.io/gorm"
)

func UpdateRanking(tx *gorm.DB, contestID uint) (bool, string) {
	if !config.Env.Redis.On {
		return false, "RedisOff"
	}
	var teams []model.Team
	res := tx.Model(&model.Team{}).Where("contest_id = ? AND banned = ?", contestID, false).
		Order("score DESC, last ASC").Find(&teams)
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

func GetRanking(contestID uint, args ...int) ([]model.Team, bool, string) {
	limit, offset := -1, 0
	if len(args) > 0 {
		if args[0] == 0 {
			args[0] = -1
		}
		limit = args[0]
	}
	if len(args) > 1 {
		offset = args[1]
	}
	if teams, err := redis.GetCachedRanking(contestID, int64(limit), int64(offset)); err == nil && teams != nil {
		return teams, true, "Success"
	}
	var teams []model.Team
	res := DB.Model(&model.Team{}).Where("contest_id = ? AND banned = ?", contestID, false).
		Order("score DESC, last ASC").Find(&teams)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get teams: %v", res.Error)
		return nil, false, "GetTeamError"
	}
	go UpdateRanking(DB, contestID)
	limit, offset = utils.TidyPaginate(len(teams), limit, offset)
	return teams[limit:offset], true, "Success"
}
