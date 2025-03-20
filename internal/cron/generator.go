package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

// PrepareGenerator 预开动态题目生成器, 后续生成附件时直接附加执行
func PrepareGenerator(c *cron.Cron) {
	log.Logger.Debug("Prepare generator")
	function := executionTime("PrepareGenerator", func() {
		var (
			ok       bool
			msg      string
			contests []model.Contest
			usages   []model.Usage
		)
		contests, _, ok, msg = db.GetContests(db.DB, 0, 0, false)
		if !ok {
			log.Logger.Warningf("Failed to get contests %s", msg)
			return
		}
		for _, contest := range contests {
			if contest.IsRunning() {
				usages, ok, msg = db.GetUsageByContestID(db.DB, contest.ID, false)
				if !ok {
					log.Logger.Warningf("Failed to get usages %s", msg)
					continue
				}
				for _, usage := range usages {
					if usage.Type == model.Dynamic {
						_, ok, msg = k8s.StartGenerator(usage)
						if !ok {
							log.Logger.Warningf("Failed to start generator %s", msg)
						}
					}
				}
			}
		}
	})
	function()
	c.Schedule(cron.Every(30*time.Minute), cron.FuncJob(function))
}

// CloseGenerator 关闭超时的动态题目生成器, 释放部分资源
func CloseGenerator(c *cron.Cron) {
	function := executionTime("CloseGenerator", func() {
		log.Logger.Debug("Close timeout generator")
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
