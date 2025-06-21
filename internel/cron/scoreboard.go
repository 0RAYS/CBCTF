package cron

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/robfig/cron/v3"
	"math"
	"time"
)

// UpdateTeamRanking 依据数据库, 更新 model.Team 的分数和排名
func UpdateTeamRanking(c *cron.Cron) {
	function := exec("UpdateTeamRanking", func() {
		repo := db.InitContestRepo(db.DB)
		contests, _, ok, _ := repo.List(-1, -1)
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
	function := exec("UpdateUserRanking", func() {
		contestRepo := db.InitContestRepo(db.DB)
		contests, _, ok, _ := contestRepo.List(-1, -1, db.GetOptions{
			Selects: []string{"id", "start", "duration"},
			Preloads: map[string]db.GetOptions{
				"Users": {
					Selects: []string{"id"},
				},
			},
		})
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
				submissions, _, ok, _ := submissionRepo.List(-1, -1, db.GetOptions{
					Conditions: map[string]any{
						"user_id": user.ID,
						"solved":  true,
					},
					Selects: []string{"id", "contest_flag_id", "team_id"},
				})
				if !ok {
					continue
				}
				var solved int64 = 0
				var score float64 = 0
				contestFlagRepo := db.InitContestFlagRepo(db.DB)
				for _, submission := range submissions {
					contestFlag, ok, _ := contestFlagRepo.GetByID(submission.ContestFlagID)
					if !ok {
						continue
					}
					solved++
					var rate float64
					bloodTeam, _, _ := submissionRepo.GetBloodTeam(submission.ContestFlagID)
					for i, teamID := range bloodTeam {
						if teamID == submission.TeamID {
							switch i {
							case 0:
								rate = model.FirstBloodRate
							case 1:
								rate = model.SecondBloodRate
							case 2:
								rate = model.ThirdBloodRate
							}
						}
						if rate > 0 {
							break
						}
					}
					score += contestFlag.CurrentScore + contestFlag.Score*rate
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
