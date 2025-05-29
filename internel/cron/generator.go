package cron

import (
	"CBCTF/internel/k8s"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"github.com/robfig/cron/v3"
	"time"
)

// PrepareGenerator 关闭超时的动态题目生成器, 释放部分资源
func PrepareGenerator(c *cron.Cron) {
	function := exec("ResetGenerator", func() {
		contestRepo := db.InitContestRepo(db.DB)
		contests, _, ok, _ := contestRepo.GetAll(-1, -1, false)
		if !ok {
			return
		}
		for _, contest := range contests {
			if !contest.IsRunning() {
				continue
			}
			repo := db.InitUsageRepo(db.DB)
			usages, _, ok, _ := repo.GetAll(contest.ID, -1, -1, false, "Challenge")
			if !ok {
				continue
			}
			for _, usage := range usages {
				if usage.Challenge.Type == model.DynamicChallenge {
					go func(usage model.Usage) {
						if _, ok, _ = k8s.StartGenerator(usage); !ok {
							k8s.StopGenerator(usage)
						}
					}(usage)
				}
			}
		}
	})
	function()
	c.Schedule(cron.Every(30*time.Minute), cron.FuncJob(function))
}
