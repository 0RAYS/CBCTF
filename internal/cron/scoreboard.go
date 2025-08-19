package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"math"
	"slices"
	"time"

	"github.com/robfig/cron/v3"
)

// updateTeamRanking 全量更新 model.Team 的分数和排名
// TODO 比赛结束时可能位于空档期，导致最终分数核算出现问题
func updateTeamRanking(c *cron.Cron) {
	function := exec("UpdateTeamRanking", func() {
		repo := db.InitContestRepo(db.DB)
		contests, _, ok, _ := repo.List(-1, -1, db.GetOptions{
			Selects:    []string{"id", "start", "duration"},
			Conditions: map[string]any{"hidden": false},
		})
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

// updateUserRanking 全量更新 model.User 的分数和排名
func updateUserRanking(c *cron.Cron) {
	function := exec("UpdateUserRanking", func() {
		userRepo := db.InitUserRepo(db.DB)
		users, _, ok, _ := userRepo.List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"banned": false},
			Selects:    []string{"id"},
			Preloads: map[string]db.GetOptions{
				"Submissions": {
					Conditions: map[string]any{"solved": true},
					Selects:    []string{"id", "user_id", "team_id", "contest_flag_id"},
					Preloads: map[string]db.GetOptions{
						"ContestFlag": {Selects: []string{"id", "current_score", "score"}},
					},
				},
			},
		})
		if !ok {
			return
		}
		submissionRepo := db.InitSubmissionRepo(db.DB)
		for _, user := range users {
			var solved int64 = 0
			var score float64 = 0
			for _, submission := range user.Submissions {
				solved++
				var rate float64
				bloodTeam, _, _ := submissionRepo.GetBloodTeam(submission.ContestFlagID)
				switch slices.IndexFunc(bloodTeam, func(i uint) bool {
					if i == submission.TeamID {
						return true
					}
					return false
				}) {
				case 0:
					rate = model.FirstBloodRate
				case 1:
					rate = model.SecondBloodRate
				case 2:
					rate = model.ThirdBloodRate
				}
				score += submission.ContestFlag.CurrentScore + submission.ContestFlag.Score*rate
			}
			score = math.Trunc(score*100) / 100
			userRepo.Update(user.ID, db.UpdateUserOptions{
				Score:  &score,
				Solved: &solved,
			})
		}
		service.UpdateUserRanking(db.DB)
	})
	function()
	c.Schedule(cron.Every(3*time.Hour), cron.FuncJob(function))
}
