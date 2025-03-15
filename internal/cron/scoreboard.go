package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"github.com/robfig/cron/v3"
	"time"
)

// UpdateTeamRanking 依据数据库, 更新 model.Team 的分数和排名
func UpdateTeamRanking(c *cron.Cron) {
	function := func() {
		log.Logger.Debug("Update global ranking")
		contests, _, ok, _ := db.GetContests(db.DB, -1, -1, false, false)
		if !ok {
			return
		}
		for _, contest := range contests {
			if !contest.IsRunning() {
				continue
			}
			go func() {
				log.Logger.Debugf("Start contest %d team ranking", contest.ID)
				db.UpdateTeamRanking(db.DB, contest.ID)
				teams, _ := redis.GetTeamRanking(contest.ID, 0, -1)
				for rank, team := range teams {
					if team.Rank == rank+1 {
						continue
					}
					tx := db.DB.Begin()
					if ok, _ := db.UpdateTeam(tx, team.ID, map[string]interface{}{"rank": rank + 1}); !ok {
						tx.Rollback()
						continue
					}
					tx.Commit()
				}
			}()
		}
	}
	function()
	c.Schedule(cron.Every(5*time.Minute), cron.FuncJob(function))
}

// UpdateUserRanking 依据数据库, 更新 model.User 的分数和排名
func UpdateUserRanking(c *cron.Cron) {
	function := func() {
		var (
			contests    []model.Contest
			users       []*model.User
			submissions []model.Submission
			ok          bool
		)
		log.Logger.Debug("Update user ranking")
		// 嵌套预加载, 需要使用 model.User.Teams 字段
		contests, _, ok, _ = db.GetContests(db.DB, -1, -1, false, true, true)
		if !ok {
			return
		}
		for _, contest := range contests {
			users = contest.Users
			for _, user := range users {
				submissions, _, ok, _ = db.GetSubmissions(db.DB, -1, -1, "user_id", user.ID)
				if !ok {
					continue
				}
				data := map[string]interface{}{
					"solved": 0,
					"score":  0.0,
				}
				for _, submission := range submissions {
					if submission.Solved {
						usage, ok, _ := db.GetUsageBy2ID(db.DB, submission.ContestID, submission.ChallengeID)
						if ok && !usage.Hidden {
							rate := func() float64 {
								for _, team := range user.Teams {
									if team.ContestID == contest.ID {
										rate, _ := usage.CalcBlood(team.ID)
										return rate
									}
								}
								return 0
							}()
							data["solved"] = data["solved"].(int) + 1
							data["score"] = data["score"].(float64) + usage.CurrentScore + usage.Score*rate
						}
					}
				}
				tx := db.DB.Begin()
				if ok, _ = db.UpdateUser(tx, user.ID, data); !ok {
					tx.Rollback()
					continue
				}
				tx.Commit()
			}
		}
		db.UpdateUserRanking(db.DB)
	}
	function()
	c.Schedule(cron.Every(12*time.Hour), cron.FuncJob(function))
}
