package cron

import (
	"CBCTF/internel/k8s"
	"CBCTF/internel/log"
	db "CBCTF/internel/repo"
	"context"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

// CloseContainers 关闭并删除数据库中存在记录的超时 Containers
func CloseContainers(c *cron.Cron) {
	function := executionTime("CloseContainers", func() {
		log.Logger.Debug("Close timeout Containers")
		repo := db.InitContainerRepo(db.DB)
		containers, _, ok, msg := repo.GetAll(-1, -1, false)
		if !ok {
			log.Logger.Warningf("Failed to get Containers %s", msg)
			return
		}
		idL := make([]uint, 0)
		for _, container := range containers {
			if container.Start.Add(container.Duration).Before(time.Now()) {
				idL = append(idL, container.ID)
				k8s.StopContainer(container)
			}
		}
		repo.Delete(idL...)
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
				repo := db.InitContainerRepo(db.DB)
				if _, ok, _ := repo.GetByName("pod", pod.Name, false, false); !ok {
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
