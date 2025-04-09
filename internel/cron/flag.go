package cron

import (
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

// UpdateFlagScore 依据数据库, 更新 model.Flag 的分数和解题人数
// 正常情况下该定时任务无意义, 每次有新解出时即更新 current_score 和 solvers
// 当 submissions 且 model.Submission.Solved == true 时的数据减少 (例如: 用户注销 / 队伍解散 引发的数据删除), 该函数才有意义
func UpdateFlagScore(c *cron.Cron) {
	function := executionTime("UpdateFlagScore", func() {
		contestRepo := db.InitContestRepo(db.DB)
		contests, _, ok, _ := contestRepo.GetAll(-1, -1, false)
		if !ok {
			return
		}
		for _, contest := range contests {
			if !contest.IsRunning() {
				continue
			}
			usageRepo := db.InitUsageRepo(db.DB)
			usages, _, ok, _ := usageRepo.GetAll(contest.ID, -1, -1, false, "Flags")
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
						if solvers != flag.Solvers || currentScore != flag.CurrentScore {
							tx := db.DB.Begin()
							if ok, _ = flagRepo.Update(flag.ID, db.UpdateFlagOptions{
								CurrentScore: &currentScore,
								Solvers:      &solvers,
							}); !ok {
								tx.Rollback()
								return
							}
							tx.Commit()
						}
					}()
					mu.(*sync.Mutex).Unlock()
				}
			}
		}
	})
	function()
	c.Schedule(cron.Every(time.Hour), cron.FuncJob(function))
}
