package cron

import (
	"CBCTF/internel/log"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

func UpdateUsageScore(c *cron.Cron) {
	function := executionTime("UpdateUsageScore", func() {
		log.Logger.Debug("Update usage score")
		contestRepo := db.InitContestRepo(db.DB)
		contests, _, ok, _ := contestRepo.GetAll(-1, -1, false, true)
		if !ok {
			return
		}
		for _, contest := range contests {
			if !contest.IsRunning() {
				continue
			}
			usageRepo := db.InitUsageRepo(db.DB)
			usages, _, ok, _ := usageRepo.GetAll(contest.ID, -1, -1, false, true)
			if !ok {
				return
			}
			for _, usage := range usages {
				for _, flag := range usage.Flags {
					mu, _ := service.SolvedMutex.LoadOrStore(flag.ID, &sync.Mutex{})
					mu.(*sync.Mutex).Lock()
					func() {
						flagRepo := db.InitFlagRepo(db.DB)
						solvers, currentScore, ok, _ := service.CalcSolversAndScore(db.DB, flag)
						if !ok {
							return
						}
						tx := db.DB.Begin()
						if ok, _ = flagRepo.Update(flag.ID, db.UpdateFlagOptions{
							CurrentScore: &currentScore,
							Solvers:      &solvers,
						}); !ok {
							tx.Rollback()
							return
						}
						tx.Commit()
					}()
					mu.(*sync.Mutex).Unlock()
				}
			}
		}
	})
	function()
	c.Schedule(cron.Every(5*time.Minute), cron.FuncJob(function))
}
