package cron

import (
	"CBCTF/internel/config"
	"CBCTF/internel/k8s"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"context"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

// PrepareGenerator 关闭超时的动态题目生成器, 释放部分资源
func PrepareGenerator(c *cron.Cron) {
	function := exec("ResetGenerator", func() {
		contests, _, ok, _ := db.InitContestRepo(db.DB).List(-1, -1)
		if !ok {
			return
		}
		contestChallengeRepo := db.InitContestChallengeRepo(db.DB)
		for _, contest := range contests {
			if contest.IsOver() {
				continue
			}
			contestChallengeL, _, ok, _ := contestChallengeRepo.ListWithConditions(-1, -1, db.GetOptions{
				{Key: "contest_id", Value: contest.ID, Op: "and"},
				{Key: "hidden", Value: false, Op: "and"},
				{Key: "type", Value: model.DynamicChallengeType, Op: "and"},
			}, false, "Challenge")
			if !ok {
				continue
			}
			for _, contestChallenge := range contestChallengeL {
				for i := 0; i < len(config.Env.K8S.Nodes)*2-len(k8s.GeneratorMap[contestChallenge.ID]); i++ {
					go k8s.StartGenerator(contestChallenge)
				}
			}
		}
	})
	function()
	c.Schedule(cron.Every(30*time.Minute), cron.FuncJob(function))
}

func StopUnCtrlGenerator(c *cron.Cron) {
	function := exec("StopUnCtrlGenerator", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		pods, ok, msg := k8s.GetPods(ctx)
		cancel()
		if !ok {
			log.Logger.Warningf("Failed to get generators %s", msg)
			return
		}
		generators := make(map[string]*k8s.Generator)
		names := make([]string, 0)
		for _, v := range k8s.GeneratorMap {
			for _, generator := range v {
				generators[generator.Pod.Name] = generator
				names = append(names, generator.Pod.Name)
			}
		}
		for _, pod := range pods.Items {
			if strings.HasPrefix(pod.Name, "gen") && !utils.In(pod.Name, names) {
				ctx, cancel = context.WithTimeout(context.Background(), 1*time.Minute)
				_, _ = k8s.DeletePod(ctx, pod.Name)
				cancel()
			}
		}
	})
	function()
	c.Schedule(cron.Every(1*time.Hour), cron.FuncJob(function))
}
