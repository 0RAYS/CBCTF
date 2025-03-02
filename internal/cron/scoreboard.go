package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/redis"
	"github.com/robfig/cron/v3"
	"time"
)

func UpdateGlobalRanking(c *cron.Cron) {
	function := func() {
		contests, _, ok, _ := db.GetContests(db.DB, -1, -1, false, false)
		if !ok {
			return
		}
		for _, contest := range contests {
			if !contest.IsRunning() {
				continue
			}
			go db.UpdateRanking(db.DB, contest.ID)
		}
	}
	function()
	c.Schedule(cron.Every(1*time.Minute), cron.FuncJob(function))
}

func UpdateTeamRank(c *cron.Cron) {
	function := func() {
		contests, _, ok, _ := db.GetContests(db.DB, -1, -1, false, false)
		if !ok {
			return
		}
		for _, contest := range contests {
			if !contest.IsRunning() {
				continue
			}
			teams, err := redis.GetCachedRanking(contest.ID, 0, -1)
			if err != nil {
				continue
			}
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
		}
	}
	function()
	c.Schedule(cron.Every(1*time.Minute), cron.FuncJob(function))
}
