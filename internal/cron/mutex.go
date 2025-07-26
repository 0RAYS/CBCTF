package cron

import (
	"CBCTF/internal/model"
	db "CBCTF/internal/repo"
	"CBCTF/internal/service"
	"github.com/robfig/cron/v3"
	"time"
)

// ClearContestChallengeMutex 定时任务清理flag提交锁 service.SolvedMutex
func ClearContestChallengeMutex(c *cron.Cron) {
	function := exec("ClearSubmissionMutex", func() {
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
	})
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(function))
}
