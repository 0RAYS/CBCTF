package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"github.com/robfig/cron/v3"
	"time"
)

// ClearUsageMutex 定时任务清理flag提交锁 db.SolvedMutex
func ClearUsageMutex(c *cron.Cron) {
	function := executionTime("ClearSubmissionMutex", func() {
		log.Logger.Debug("Clear submission mutex")
		contests := make(map[uint]model.Contest)
		db.SolvedMutex.Range(func(k, v interface{}) bool {
			usage, ok, _ := db.GetUsageByID(db.DB, k.(uint))
			if !ok {
				db.SolvedMutex.Delete(k)
				return true
			}
			if contest, ok := contests[usage.ContestID]; !ok {
				contest, ok, _ = db.GetContestByID(db.DB, usage.ContestID)
				if !ok {
					db.SolvedMutex.Delete(k)
					return true
				}
				contests[usage.ContestID] = contest
				if !contest.IsRunning() {
					db.SolvedMutex.Delete(k)
				}
			}
			return true
		})
	})
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(function))
}
