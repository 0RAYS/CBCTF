package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"time"

	"github.com/robfig/cron/v3"
)

// ClearContestChallengeMutex 定时任务清理flag提交锁 service.SolvedMutex
func ClearContestChallengeMutex(c *cron.Cron) {
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(exec("ClearSubmissionMutex", func() {
		contests := make(map[uint]model.Contest)
		contestRepo := db.InitContestRepo(db.DB)
		contestFlagRepo := db.InitContestFlagRepo(db.DB)
		service.SolvedMutex.Range(func(k, v any) bool {
			contestFlag, ok, _ := contestFlagRepo.GetByID(k.(uint))
			if !ok {
				service.SolvedMutex.Delete(k)
				return true
			}
			if contest, ok := contests[contestFlag.ContestID]; !ok {
				contest, ok, _ = contestRepo.GetByID(contestFlag.ContestID)
				if !ok {
					service.SolvedMutex.Delete(k)
					return true
				}
				contests[contestFlag.ContestID] = contest
				if !contest.IsRunning() {
					service.SolvedMutex.Delete(k)
				}
			}
			return true
		})
	})))
}
