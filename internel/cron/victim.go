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

func CloseTimeoutVictims(c *cron.Cron) {
	function := exec("CloseTimeoutVictims", func() {
		repo := db.InitVictimRepo(db.DB)
		victims, _, ok, _ := repo.List(-1, -1, db.GetOptions{
			Selects: []string{"id", "start", "duration", "team_id", "contest_challenge_id"},
			Preloads: map[string]db.GetOptions{
				"Pods": {
					Selects: []string{"id", "name", "victim_id"},
				},
			},
		})
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

func CloseUnCtrlVictims(c *cron.Cron) {
	function := exec("CloseUnCtrlVictims", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		pods, ok, msg := k8s.GetPods(ctx)
		cancel()
		if !ok {
			log.Logger.Warningf("Failed to get victims %s", msg)
			return
		}
		podRepo := db.InitPodRepo(db.DB)
		for _, pod := range pods.Items {
			if strings.HasPrefix(pod.Name, k8s.VictimPodTag) && time.Now().Sub(pod.CreationTimestamp.Time) > 4*time.Hour {
				_, ok, _ = podRepo.Get(db.GetOptions{
					Conditions: map[string]any{"name": pod.Name},
					Selects:    []string{"id"},
				})
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
