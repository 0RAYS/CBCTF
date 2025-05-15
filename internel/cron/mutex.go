package cron

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/robfig/cron/v3"
	"time"
)

// ClearUsageMutex 定时任务清理flag提交锁 service.SolvedMutex
func ClearUsageMutex(c *cron.Cron) {
	function := exec("ClearSubmissionMutex", func() {
		contests := make(map[uint]model.Contest)
		contestRepo := db.InitContestRepo(db.DB)
		flagRepo := db.InitFlagRepo(db.DB)
		service.SolvedMutex.Range(func(k, v any) bool {
			flag, ok, _ := flagRepo.GetByID(k.(uint))
			if !ok {
				service.SolvedMutex.Delete(k)
				return true
			}
			if contest, ok := contests[flag.ContestID]; !ok {
				contest, ok, _ = contestRepo.GetByID(flag.ContestID)
				if !ok {
					service.SolvedMutex.Delete(k)
					return true
				}
				contests[flag.ContestID] = contest
				if !contest.IsRunning() {
					service.SolvedMutex.Delete(k)
				}
			}
			return true
		})
	})
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(function))
}
