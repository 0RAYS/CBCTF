package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"

	"gorm.io/gorm"
)

func UpdateTeamRanking(tx *gorm.DB, contest model.Contest, limit, offset int) ([]model.Team, int64, model.RetVal) {
	var (
		repo          = db.InitTeamRepo(tx)
		teams, _, ret = repo.List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"contest_id": contest.ID, "banned": false},
		})
		score float64
	)
	if !ret.OK {
		return nil, 0, ret
	}
	for i, team := range teams {
		score, ret = CalcTeamScore(tx, team, contest.Blood)
		if !ret.OK {
			continue
		}
		teams[i].Score = score
	}
	if ret = redis.UpdateTeamRanking(contest.ID, teams); !ret.OK {
		return nil, 0, ret
	}
	teams, count, ret := GetTeamRanking(tx, contest, limit, offset)
	if !ret.OK {
		return nil, 0, ret
	}
	for i, team := range teams {
		teams[i].Rank = i + 1
		repo.Update(team.ID, db.UpdateTeamOptions{Score: &team.Score, Rank: new(i + 1)})
	}
	return teams, count, model.SuccessRetVal()
}

func GetTeamRanking(tx *gorm.DB, contest model.Contest, limit, offset int) ([]model.Team, int64, model.RetVal) {
	count, ret := db.InitTeamRepo(tx).Count(db.CountOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "banned": false},
	})
	if !ret.OK {
		return nil, 0, ret
	}
	start, end := utils.TidyPaginate(int(count), limit, offset)
	if end-start <= 0 {
		return nil, count, model.SuccessRetVal()
	}
	teams, ret := redis.GetTeamRanking(contest.ID, int64(start), int64(end-1))
	if !ret.OK || (end-start > 0 && len(teams) == 0 && count > 0) {
		return UpdateTeamRanking(tx, contest, limit, offset)
	}
	return teams, count, model.SuccessRetVal()
}

func UpdateUserRanking(tx *gorm.DB, limit, offset int) ([]model.User, int64, model.RetVal) {
	users, _, ret := db.InitUserRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"banned": false},
	})
	if !ret.OK {
		return nil, 0, ret
	}
	if ret = redis.UpdateUserRanking(users); !ret.OK {
		return nil, 0, ret
	}
	return GetUserRanking(tx, limit, offset)
}

func GetUserRanking(tx *gorm.DB, limit, offset int) ([]model.User, int64, model.RetVal) {
	count, ret := db.InitUserRepo(tx).Count(db.CountOptions{
		Conditions: map[string]any{"banned": false},
	})
	if !ret.OK {
		return nil, count, ret
	}
	start, end := utils.TidyPaginate(int(count), limit, offset)
	if end-start <= 0 {
		return nil, count, model.SuccessRetVal()
	}
	users, ret := redis.GetUserRanking(int64(start), int64(end-1))
	if !ret.OK || (end-start > 0 && len(users) == 0 && count > 0) {
		return UpdateUserRanking(tx, limit, offset)
	}
	return users, count, model.SuccessRetVal()
}
