package cron

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/robfig/cron/v3"
	"time"
)

// ClearUsageMutex 定时任务清理flag提交锁 service.SolvedMutex
func ClearUsageMutex(c *cron.Cron) {
	function := executionTime("ClearSubmissionMutex", func() {
		log.Logger.Debug("Clear submission mutex")
		contests := make(map[uint]model.Contest)
		contestRepo := db.InitContestRepo(db.DB)
		usageRepo := db.InitUsageRepo(db.DB)
		service.SolvedMutex.Range(func(k, v interface{}) bool {
			usage, ok, _ := usageRepo.GetByID(k.(uint), false)
			if !ok {
				service.SolvedMutex.Delete(k)
				return true
			}
			if contest, ok := contests[usage.ContestID]; !ok {
				contest, ok, _ = contestRepo.GetByID(usage.ContestID, false)
				if !ok {
					service.SolvedMutex.Delete(k)
					return true
				}
				contests[usage.ContestID] = contest
				if !contest.IsRunning() {
					service.SolvedMutex.Delete(k)
				}
			}
			return true
		})
	})
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(function))
}
