package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"

	"gorm.io/gorm"
)

func UpdateTeamRanking(tx *gorm.DB, contest model.Contest, limit, offset int) ([]model.Team, int64, bool, string) {
	var (
		repo              = db.InitTeamRepo(tx)
		teams, _, ok, msg = repo.List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"contest_id": contest.ID, "banned": false},
		})
		score float64
	)
	if !ok {
		return nil, 0, false, msg
	}
	for i, team := range teams {
		score, ok, msg = CalcTeamScore(tx, team)
		if !ok {
			continue
		}
		teams[i].Score = score
	}
	if err := redis.UpdateTeamRanking(contest.ID, teams); err != nil {
		log.Logger.Warningf("Failed to update TeamRanking: %s", err)
		return nil, 0, false, i18n.UpdateRankingError
	}
	teams, count, ok, msg := GetTeamRanking(tx, contest, limit, offset)
	if !ok {
		return nil, 0, false, msg
	}
	for i, team := range teams {
		teams[i].Rank = i + 1
		repo.Update(team.ID, db.UpdateTeamOptions{Score: &team.Score, Rank: utils.Ptr(i + 1)})
	}
	return teams, count, true, i18n.Success
}

func GetTeamRanking(tx *gorm.DB, contest model.Contest, limit, offset int) ([]model.Team, int64, bool, string) {
	count, ok, msg := db.InitTeamRepo(tx).Count(db.CountOptions{
		Conditions: map[string]any{"contest_id": contest.ID, "banned": false},
	})
	if !ok {
		return nil, 0, false, msg
	}
	start, end := utils.TidyPaginate(int(count), limit, offset)
	if end-start <= 0 {
		return nil, count, true, i18n.Success
	}
	teams, err := redis.GetTeamRanking(contest.ID, int64(start), int64(end-1))
	if err != nil || (end-start > 0 && len(teams) == 0 && count > 0) {
		return UpdateTeamRanking(tx, contest, limit, offset)
	}
	return teams, count, true, i18n.Success
}

func UpdateUserRanking(tx *gorm.DB, limit, offset int) ([]model.User, int64, bool, string) {
	users, _, ok, msg := db.InitUserRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"banned": false},
	})
	if !ok {
		return nil, 0, false, msg
	}
	if err := redis.UpdateUserRanking(users); err != nil {
		log.Logger.Warningf("Failed to update UserRanking: %s", err)
		return nil, 0, false, i18n.UpdateRankingError
	}
	return GetUserRanking(tx, limit, offset)
}

func GetUserRanking(tx *gorm.DB, limit, offset int) ([]model.User, int64, bool, string) {
	count, ok, msg := db.InitUserRepo(tx).Count(db.CountOptions{
		Conditions: map[string]any{"banned": false},
	})
	if !ok {
		return nil, count, false, msg
	}
	start, end := utils.TidyPaginate(int(count), limit, offset)
	if end-start <= 0 {
		return nil, count, true, i18n.Success
	}
	users, err := redis.GetUserRanking(int64(start), int64(end-1))
	if err != nil || (end-start > 0 && len(users) == 0 && count > 0) {
		return UpdateUserRanking(tx, limit, offset)
	}
	return users, count, true, i18n.Success
}
