package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"time"

	"github.com/robfig/cron/v3"
)

// clearSubmissionMutex 定时任务清理flag提交锁 service.SolvedMutex
func clearSubmissionMutex(c *cron.Cron) {
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(exec("ClearSubmissionMutex", func() {
		contests := make(map[uint]model.Contest)
		contestRepo := db.InitContestRepo(db.DB)
		contestFlagRepo := db.InitContestFlagRepo(db.DB)
		service.SolvedMutex.Range(func(k, v any) bool {
			contestFlag, ret := contestFlagRepo.GetByID(k.(uint))
			if !ret.OK {
				service.SolvedMutex.Delete(k)
				return true
			}
			contest, ok := contests[contestFlag.ContestID]
			if !ok {
				contest, ret = contestRepo.GetByID(contestFlag.ContestID)
				if !ret.OK {
					service.SolvedMutex.Delete(k)
					return true
				}
				contests[contestFlag.ContestID] = contest
			}
			if !contest.IsRunning() {
				service.SolvedMutex.Delete(k)
			}
			return true
		})
	})))
}
