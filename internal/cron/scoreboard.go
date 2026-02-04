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
		})
		if !ret.OK {
			return
		}
		userIDs := make([]uint, len(users))
		for i, user := range users {
			userIDs[i] = user.ID
		}

		// 批量查询所有用户的已解决提交 - 优化：从嵌套 Preload 改为单次 JOIN 查询
		submissionRepo := db.InitSubmissionRepo(db.DB)
		submissions, ret := submissionRepo.GetUserSolvedSubmissions(userIDs...)
		if !ret.OK {
			return
		}

		// 构建 userID -> submissions 映射
		userSubmissionsMap := make(map[uint][]db.UserSolvedSubmission)
		contestIDSet := make(map[uint]bool)
		for _, submission := range submissions {
			userSubmissionsMap[submission.UserID] = append(userSubmissionsMap[submission.UserID], submission)
			contestIDSet[submission.ContestID] = true
		}

		// 查询涉及的 contests 以获取 blood 设置
		contestIDL := make([]uint, 0, len(contestIDSet))
		for contestID := range contestIDSet {
			contestIDL = append(contestIDL, contestID)
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
		for _, user := range users {
			userSubmissions := userSubmissionsMap[user.ID]
			var solved int64 = 0
			var score float64 = 0
			for _, submission := range userSubmissions {
				solved++
				var rate float64
				if blood[submission.ContestID] {
					bloodTeam, _ := submissionRepo.GetBloodTeamID(submission.ContestFlagID)
					switch slices.IndexFunc(bloodTeam, func(i uint) bool {
						return i == submission.TeamID
					}) {
					case 0:
						rate = model.FirstBloodRate
					case 1:
						rate = model.SecondBloodRate
					case 2:
						rate = model.ThirdBloodRate
					}
				}
				score += submission.ContestFlagCurrentScore + submission.ContestFlagScore*rate
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
