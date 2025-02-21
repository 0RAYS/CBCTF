package cron

import (
	"CBCTF/internal/db"
	"github.com/robfig/cron/v3"
	"time"
)

func UpdateRanking(c *cron.Cron) {
	function := func() {
		contests, _, ok, _ := db.GetContests(db.DB, -1, -1, false, false)
		if !ok {
			return
		}
		for _, contest := range contests {
			if !contest.IsRunning() {
				return
			}
			go db.UpdateRanking(db.DB, contest.ID)
		}
	}
	function()
	c.Schedule(cron.Every(1*time.Minute), cron.FuncJob(function))
}
