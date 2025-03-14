package cron

import (
	"CBCTF/internal/db"
	"github.com/robfig/cron/v3"
	"time"
)

func ClearUsageMutex(c *cron.Cron) {
	function := func() {
		db.UsagesMutex.Range(func(k, v interface{}) bool {
			usage, ok, _ := db.GetUsageByID(db.DB, k.(uint))
			if !ok {
				return true
			}
			contest, ok, _ := db.GetContestByID(db.DB, usage.ContestID)
			if !ok {
				return true
			}
			if !contest.IsRunning() {
				db.UsagesMutex.Delete(k)
			}
			return true
		})
	}
	function()
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(function))
}
