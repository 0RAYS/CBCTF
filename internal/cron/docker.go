package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

// CloseDockers 关闭并删除超时 dockers
func CloseDockers(c *cron.Cron) {
	function := func() {
		log.Logger.Info("Close timeout dockers")
		dockers, ok, msg := db.GetDockers(db.DB, false)
		if !ok {
			log.Logger.Warningf("Failed to get dockers %s", msg)
			return
		}
		for _, docker := range dockers {
			if docker.Start.Add(docker.Duration).Before(time.Now()) {
				// 每次删除都作为一个单独的事务, 不回滚之前的删除
				tx := db.DB.Begin()
				if ok, msg = db.DeleteDocker(tx, docker); !ok {
					tx.Rollback()
					log.Logger.Warningf("Failed to delete docker %s", msg)
					continue
				}
				tx.Commit()
			}
		}
	}
	function()
	c.Schedule(cron.Every(5*time.Minute), cron.FuncJob(function))
}

// CloseUnCtrlDockers 移除意外超时的 pod
func CloseUnCtrlDockers(c *cron.Cron) {
	function := func() {
		log.Logger.Info("Close timeout pods")
		pods, ok, msg := k8s.GetPods()
		if !ok {
			log.Logger.Warningf("Failed to get pods %s", msg)
			return
		}
		for _, pod := range pods.Items {
			if strings.Contains(pod.Name, "victim") && time.Now().Sub(pod.CreationTimestamp.Time) > 4*time.Hour {
				if _, ok, _ := db.GetDockerByPodName(db.DB, pod.Name); !ok {
					_, _ = k8s.DeletePod(pod.Name)
				}
			}
		}
	}
	function()
	c.Schedule(cron.Every(1*time.Hour), cron.FuncJob(function))
}
