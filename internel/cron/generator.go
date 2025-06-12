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
		contests, _, ok, _ := contestRepo.List(-1, -1)
		if !ok {
			return
		}
		for _, contest := range contests {
			if contest.IsOver() {
				continue
			}
			repo := db.InitContestChallengeRepo(db.DB)
			contestChallengeL, _, ok, _ := repo.ListWithConditions(-1, -1, db.GetOptions{
				{Key: "contest_id", Value: contest.ID, Op: "and"},
				{Key: "hidden", Value: false, Op: "and"},
			}, false, "Challenge")
			if !ok {
				continue
			}
			for _, contestChallenge := range contestChallengeL {
				if contestChallenge.Challenge.Type == model.DynamicChallengeType {
					go func(contestChallenge model.ContestChallenge) {
						if _, ok, _ = k8s.StartGenerator(contestChallenge); !ok {
							k8s.StopGenerator(contestChallenge)
						}
					}(contestChallenge)
				}
			}
		}
	})
	function()
	c.Schedule(cron.Every(30*time.Minute), cron.FuncJob(function))
}
