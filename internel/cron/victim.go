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

func ClearUnCtrlResource(c *cron.Cron) {
	function := exec("ClearUnCtrlResource", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		if configmaps, ok, _ := k8s.GetConfigMapList(ctx); ok {
			for _, cm := range configmaps.Items {
				for k, v := range cm.Labels {
					if k == "victim" {
						if _, ok, _ = k8s.GetPod(ctx, v); !ok {
							k8s.DeleteConfigMapListByPodName(ctx, v)
						}
					}
				}
			}
		}
		cancel()
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
		if configmaps, ok, _ := k8s.GetServiceList(ctx); ok {
			for _, cm := range configmaps.Items {
				for k, v := range cm.Labels {
					if k == "victim" {
						if _, ok, _ = k8s.GetPod(ctx, v); !ok {
							k8s.DeleteServiceListByPodName(ctx, v)
						}
					}
				}
			}
		}
		cancel()
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
		if configmaps, ok, _ := k8s.GetNetworkPolicyList(ctx); ok {
			for _, cm := range configmaps.Items {
				for k, v := range cm.Labels {
					if k == "victim" {
						if _, ok, _ = k8s.GetPod(ctx, v); !ok {
							k8s.DeleteNetworkPolicyListByPodName(ctx, v)
						}
					}
				}
			}
		}
		cancel()
	})
	function()
	c.Schedule(cron.Every(time.Hour), cron.FuncJob(function))
}

func StopTimeoutVictims(c *cron.Cron) {
	function := exec("CloseTimeoutVictims", func() {
		repo := db.InitVictimRepo(db.DB)
		victims, _, ok, _ := repo.List(-1, -1, "Pods")
		if !ok {
			return
		}
		idL := make([]uint, 0)
		for _, victim := range victims {
			if victim.Start.Add(victim.Duration).Before(time.Now()) {
				idL = append(idL, victim.ID)
				k8s.StopVictim(victim)
			}
		}
		repo.Delete(idL...)
	})
	function()
	c.Schedule(cron.Every(5*time.Minute), cron.FuncJob(function))
}

func StopUnCtrlPods(c *cron.Cron) {
	function := exec("CloseUnCtrlPods", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		pods, ok, msg := k8s.GetPods(ctx)
		cancel()
		if !ok {
			log.Logger.Warningf("Failed to get pods %s", msg)
			return
		}
		podRepo := db.InitPodRepo(db.DB)
		for _, pod := range pods.Items {
			if strings.Contains(pod.Name, "victim") && time.Now().Sub(pod.CreationTimestamp.Time) > 4*time.Hour {
				_, ok, _ = podRepo.GetWithConditions(db.GetOptions{
					{Key: "name", Value: pod.Name, Op: "and"},
				}, false)
				if !ok {
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
