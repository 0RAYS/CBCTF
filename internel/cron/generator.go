package cron

import (
	"CBCTF/internel/k8s"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"context"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

// PrepareGenerator 预开动态题目生成器, 后续生成附件时直接附加执行
func PrepareGenerator(c *cron.Cron) {
	function := exec("PrepareGenerator", func() {
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
				if usage.Challenge.Type != model.DynamicChallenge {
					continue
				}
				k8s.StartGenerator(usage)
			}
		}
	})
	function()
	c.Schedule(cron.Every(30*time.Minute), cron.FuncJob(function))
}

// CloseGenerator 关闭超时的动态题目生成器, 释放部分资源
func CloseGenerator(c *cron.Cron) {
	function := exec("CloseGenerator", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		pods, ok, msg := k8s.GetPods(ctx)
		cancel()
		if !ok {
			log.Logger.Warningf("Failed to get pods %s", msg)
			return
		}
		for _, pod := range pods.Items {
			if strings.Contains(pod.Name, "generator") && time.Now().Sub(pod.CreationTimestamp.Time) > 3*time.Hour {
				ctx, cancel = context.WithTimeout(context.Background(), 1*time.Minute)
				ok, msg = k8s.DeletePod(ctx, pod.Name)
				cancel()
				if !ok {
					log.Logger.Warningf("Failed to delete pod %s %s", pod.Name, msg)
				}
			}
		}
	})
	function()
	c.Schedule(cron.Every(1*time.Hour), cron.FuncJob(function))
}
