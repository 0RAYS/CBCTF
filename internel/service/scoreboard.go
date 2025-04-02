package service

import (
	"CBCTF/internel/model"
	"CBCTF/internel/redis"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

func UpdateTeamRanking(tx *gorm.DB, contestID uint) (bool, string) {
	var (
		repo              = db.InitTeamRepo(tx)
		teams, _, ok, msg = repo.GetAll(contestID, -1, -1, true, 0, true, false)
		score             float64
		err               error
	)
	if !ok {
		return false, msg
	}
	for _, team := range teams {
		score, ok, msg = CalcTeamScore(tx, team.ID)
		if !ok {
			continue
		}
		// 不考虑更新失败的情况, 不回滚
		repo.Update(team.ID, db.UpdateTeamOptions{Score: &score})
	}
	teams, _, ok, msg = repo.GetAll(contestID, -1, -1, true, 0, true, false)
	if !ok {
		return false, msg
	}
	if err = redis.UpdateTeamRanking(contestID, teams); err != nil {
		return false, "UpdateRankingError"
	}
	return true, "Success"
}

func GetTeamRanking(tx *gorm.DB, contestID uint, limit, offset int) ([]model.Team, int64, bool, string) {
	var (
		teams          = make([]model.Team, 0)
		repo           = db.InitTeamRepo(tx)
		count, ok, msg = repo.Count(contestID, true, false)
		err            error
	)
	if !ok {
		return teams, count, false, msg
	}
	start, end := utils.TidyPaginate(int(count), limit, offset)
	teams, err = redis.GetTeamRanking(contestID, int64(start), int64(end-1))
	if err == nil {
		return teams, count, true, "Success"
	}
	if ok, msg = UpdateTeamRanking(tx, contestID); !ok {
		return teams, count, false, msg
	}
	return GetTeamRanking(tx, contestID, limit, offset)
}

func UpdateUserRanking(tx *gorm.DB) (bool, string) {
	var (
		repo              = db.InitUserRepo(tx)
		users, _, ok, msg = repo.GetAll(-1, -1, true, 0, true, false)
		err               error
	)
	if !ok {
		return false, msg
	}
	err = redis.UpdateUserRanking(users)
	if err != nil {
		return false, "UpdateRankingError"
	}
	return true, "Success"
}

func GetUserRanking(tx *gorm.DB, limit, offset int) ([]model.User, int64, bool, string) {
	var (
		users          = make([]model.User, 0)
		repo           = db.InitUserRepo(tx)
		count, ok, msg = repo.Count(true, false)
		err            error
	)
	if !ok {
		return users, count, false, msg
	}
	start, end := utils.TidyPaginate(int(count), limit, offset)
	if users, err = redis.GetUserRanking(int64(start), int64(end-1)); err == nil {
		return users, count, true, "Success"
	}
	if ok, msg = UpdateUserRanking(tx); !ok {
		return users, count, false, msg
	}
	return GetUserRanking(tx, limit, offset)
}
