package service

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	db "CBCTF/internal/repo"
	"CBCTF/internal/utils"

	"gorm.io/gorm"
)

func UpdateTeamRanking(tx *gorm.DB, contestID uint) (bool, string) {
	var (
		repo              = db.InitTeamRepo(tx)
		teams, _, ok, msg = repo.List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"contest_id": contestID, "banned": false},
			Preloads:   map[string]db.GetOptions{"Users": {}},
		})
		score float64
		err   error
	)
	if !ok {
		return false, msg
	}
	for i, team := range teams {
		score, ok, msg = CalcTeamScore(tx, team)
		if !ok {
			continue
		}
		// 不考虑更新失败的情况, 不回滚
		ok, _ = repo.Update(team.ID, db.UpdateTeamOptions{Score: &score})
		if ok {
			teams[i].Score = score
		}
	}
	if err = redis.UpdateTeamRanking(contestID, teams); err != nil {
		log.Logger.Warningf("Failed to update TeamRanking: %s", err)
		return false, i18n.UpdateRankingError
	}
	return true, i18n.Success
}

func GetTeamRanking(tx *gorm.DB, contestID uint, limit, offset int) ([]model.Team, int64, bool, string) {
	var (
		teams          = make([]model.Team, 0)
		count, ok, msg = db.InitTeamRepo(tx).Count(db.CountOptions{
			Conditions: map[string]any{"contest_id": contestID, "banned": false},
		})
		err error
	)
	if !ok {
		return teams, count, false, msg
	}
	start, end := utils.TidyPaginate(int(count), limit, offset)
	if end-start <= 0 {
		return teams, count, true, i18n.Success
	}
	teams, err = redis.GetTeamRanking(contestID, int64(start), int64(end-1))
	if err != nil || (end-start > 0 && len(teams) == 0) {
		if ok, msg = UpdateTeamRanking(tx, contestID); !ok {
			return teams, count, false, msg
		}
		return GetTeamRanking(tx, contestID, limit, offset)
	}
	return teams, count, true, i18n.Success
}

func UpdateUserRanking(tx *gorm.DB) (bool, string) {
	var (
		users, _, ok, msg = db.InitUserRepo(tx).List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"banned": false},
		})
		err error
	)
	if !ok {
		return false, msg
	}
	if err = redis.UpdateUserRanking(users); err != nil {
		log.Logger.Warningf("Failed to update UserRanking: %s", err)
		return false, i18n.UpdateRankingError
	}
	return true, i18n.Success
}

func GetUserRanking(tx *gorm.DB, limit, offset int) ([]model.User, int64, bool, string) {
	var (
		users          = make([]model.User, 0)
		count, ok, msg = db.InitUserRepo(tx).Count(db.CountOptions{
			Conditions: map[string]any{"banned": false},
		})
		err error
	)
	if !ok {
		return users, count, false, msg
	}
	start, end := utils.TidyPaginate(int(count), limit, offset)
	if end-start <= 0 {
		return users, count, true, i18n.Success
	}
	users, err = redis.GetUserRanking(int64(start), int64(end-1))
	if err != nil || (end-start > 0 && len(users) == 0) {
		if ok, msg = UpdateUserRanking(tx); !ok {
			return users, count, false, msg
		}
		return GetUserRanking(tx, limit, offset)
	}
	return users, count, true, i18n.Success
}
