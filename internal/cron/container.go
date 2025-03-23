package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"context"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

// CloseContainers 关闭并删除数据库中存在记录的超时 Containers
func CloseContainers(c *cron.Cron) {
	function := executionTime("CloseContainers", func() {
		log.Logger.Debug("Close timeout Containers")
		containers, ok, msg := db.GetContainers(db.DB, false)
		if !ok {
			log.Logger.Warningf("Failed to get Containers %s", msg)
			return
		}
		for _, container := range containers {
			if container.Start.Add(container.Duration).Before(time.Now()) {
				// 每次删除都作为一个单独的事务, 不回滚之前的删除
				tx := db.DB.Begin()
				if ok, msg = db.DeleteContainer(tx, container); !ok {
					tx.Rollback()
					log.Logger.Warningf("Failed to delete Container %s", msg)
					continue
				}
				tx.Commit()
			}
		}
	})
	function()
	c.Schedule(cron.Every(5*time.Minute), cron.FuncJob(function))
}

// CloseUnCtrlContainers 移除数据库中无记录的超时 Containers
func CloseUnCtrlContainers(c *cron.Cron) {
	function := executionTime("CloseUnCtrlContainers", func() {
		log.Logger.Debug("Close timeout pods")
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		pods, ok, msg := k8s.GetPods(ctx)
		cancel()
		if !ok {
			log.Logger.Warningf("Failed to get pods %s", msg)
			return
		}
		for _, pod := range pods.Items {
			if strings.Contains(pod.Name, "victim") && time.Now().Sub(pod.CreationTimestamp.Time) > 4*time.Hour {
				if _, ok, _ := db.GetContainerByPodName(db.DB, pod.Name); !ok {
					ctx, cancel = context.WithTimeout(context.Background(), 1*time.Minute)
					_, _ = k8s.DeletePod(ctx, pod.Name)
					cancel()
				}
			}
		}
	})
	function()
	c.Schedule(cron.Every(1*time.Hour), cron.FuncJob(function))
}
