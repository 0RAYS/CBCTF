package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"github.com/robfig/cron/v3"
	"time"
)

// ClearUsageMutex 定时任务清理flag提交锁 db.UsagesMutex
func ClearUsageMutex(c *cron.Cron) {
	function := func() {
		var contests map[uint]model.Contest
		db.UsagesMutex.Range(func(k, v interface{}) bool {
			usage, ok, _ := db.GetUsageByID(db.DB, k.(uint))
			if !ok {
				return true
			}
			if contest, ok := contests[usage.ContestID]; !ok {
				contest, ok, _ = db.GetContestByID(db.DB, usage.ContestID)
				if !ok {
					return true
				}
				contests[usage.ContestID] = contest
				if !contest.IsRunning() {
					db.UsagesMutex.Delete(k)
				}
			}
			return true
		})
	}
	function()
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(function))
}
