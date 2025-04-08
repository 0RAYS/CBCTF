package cron

import (
	"CBCTF/internel/log"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/robfig/cron/v3"
	"math"
	"time"
)

// UpdateTeamRanking 依据数据库, 更新 model.Team 的分数和排名
func UpdateTeamRanking(c *cron.Cron) {
	function := executionTime("UpdateTeamRanking", func() {
		log.Logger.Debug("Update global ranking")
		repo := db.InitContestRepo(db.DB)
		contests, _, ok, _ := repo.GetAll(-1, -1, false)
		if !ok {
			return
		}
		for _, contest := range contests {
			if !contest.IsRunning() {
				continue
			}
			service.UpdateTeamRanking(db.DB, contest.ID)
			teams, _, ok, _ := service.GetTeamRanking(db.DB, contest.ID, -1, -1)
			if !ok {
				continue
			}
			teamRepo := db.InitTeamRepo(db.DB)
			for i, team := range teams {
				rank := i + 1
				teamRepo.Update(team.ID, db.UpdateTeamOptions{Rank: &rank})
			}
		}
	})
	function()
	c.Schedule(cron.Every(5*time.Minute), cron.FuncJob(function))
}

// UpdateUserRanking 依据数据库, 更新 model.User 的分数和排名
func UpdateUserRanking(c *cron.Cron) {
	function := executionTime("UpdateUserRanking", func() {
		contestRepo := db.InitContestRepo(db.DB)
		contests, _, ok, _ := contestRepo.GetAll(-1, -1, false, "Users")
		if !ok {
			return
		}
		for _, contest := range contests {
			if !contest.IsRunning() {
				continue
			}
			submissionRepo := db.InitSubmissionRepo(db.DB)
			userRepo := db.InitUserRepo(db.DB)
			for _, user := range contest.Users {
				submissions, _, ok, _ := submissionRepo.GetAllByKeyID("user_id", user.ID, -1, -1, true)
				if !ok {
					continue
				}
				var solved int64 = 0
				var score float64 = 0
				flagRepo := db.InitFlagRepo(db.DB)
				for _, submission := range submissions {
					flag, ok, _ := flagRepo.GetByID(submission.FlagID)
					if !ok {
						continue
					}
					solved++
					rate, _ := flag.CalcBlood(submission.TeamID)
					score += flag.CurrentScore + flag.Score*rate
				}
				score = math.Trunc(score*100) / 100
				userRepo.Update(user.ID, db.UpdateUserOptions{
					Score:  &score,
					Solved: &solved,
				})
			}
		}
		service.UpdateUserRanking(db.DB)
	})
	function()
	c.Schedule(cron.Every(3*time.Hour), cron.FuncJob(function))
}
