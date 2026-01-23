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
func updateTeamRanking(c *cron.Cron) {
	function := exec("UpdateTeamRanking", func() {
		repo := db.InitContestRepo(db.DB)
		contests, _, ret := repo.List(-1, -1, db.GetOptions{
			Selects:    []string{"id", "start", "duration", "blood"},
			Conditions: map[string]any{"hidden": false},
		})
		if !ret.OK {
			return
		}
		for _, contest := range contests {
			if time.Now().Sub(contest.Start.Add(contest.Duration)) > time.Minute*10 {
				continue
			}
			service.UpdateTeamRanking(db.DB, contest, -1, -1)
		}
	})
	function()
	c.Schedule(cron.Every(5*time.Minute), cron.FuncJob(function))
}

// updateUserRanking 全量更新 model.User 的分数和排名
func updateUserRanking(c *cron.Cron) {
	function := exec("UpdateUserRanking", func() {
		userRepo := db.InitUserRepo(db.DB)
		users, _, ret := userRepo.List(-1, -1, db.GetOptions{
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
		if !ret.OK {
			return
		}
		contestIDL := make([]uint, 0)
		for _, user := range users {
			for _, submission := range user.Submissions {
				if !slices.Contains(contestIDL, submission.ContestID) {
					contestIDL = append(contestIDL, submission.ContestID)
				}
			}
		}
		contests, _, ret := db.InitContestRepo(db.DB).List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"id": contestIDL},
			Selects:    []string{"id", "blood"},
		})
		if !ret.OK {
			return
		}
		blood := make(map[uint]bool)
		for _, contest := range contests {
			blood[contest.ID] = contest.Blood
		}
		submissionRepo := db.InitSubmissionRepo(db.DB)
		for _, user := range users {
			var solved int64 = 0
			var score float64 = 0
			for _, submission := range user.Submissions {
				solved++
				var rate float64
				if a, _ := blood[submission.ContestID]; a {
					bloodTeam, _ := submissionRepo.GetBloodTeam(submission.ContestFlagID)
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
				}
				score += submission.ContestFlag.CurrentScore + submission.ContestFlag.Score*rate
			}
			score = math.Trunc(score*100) / 100
			userRepo.Update(user.ID, db.UpdateUserOptions{
				Score:  &score,
				Solved: &solved,
			})
		}
		service.UpdateUserRanking(db.DB, -1, -1)
	})
	function()
	c.Schedule(cron.Every(3*time.Hour), cron.FuncJob(function))
}
