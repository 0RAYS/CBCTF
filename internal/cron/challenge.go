package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

func UpdateUsageScore(c *cron.Cron) {
	function := executionTime("UpdateUsageScore", func() {
		log.Logger.Debug("Update usage score")
		contests, _, ok, _ := db.GetContests(db.DB, -1, -1, false, true, true)
		if !ok {
			return
		}
		for _, contest := range contests {
			if !contest.IsRunning() {
				continue
			}
			usages, _, _ := db.GetUsageByContestID(db.DB, contest.ID, true)
			for _, usage := range usages {
				mu, _ := db.SolvedMutex.LoadOrStore(usage.ID, &sync.Mutex{})
				mu.(*sync.Mutex).Lock()
				func(usage model.Usage) {
					solvers, currentScore, ok, _ := db.CalcNewUsage(db.DB, usage)
					if !ok {
						return
					}
					tx := db.DB.Begin()
					if ok, _ = db.UpdateUsage(tx, usage.ID, map[string]interface{}{
						"solvers":       solvers,
						"current_score": currentScore,
					}); !ok {
						tx.Rollback()
						return
					}
					tx.Commit()
				}(usage)
				mu.(*sync.Mutex).Unlock()
			}
		}
	})
	function()
	c.Schedule(cron.Every(5*time.Minute), cron.FuncJob(function))
}
