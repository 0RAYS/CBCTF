package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/service"
	"context"
	"slices"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
)

// closeTimeoutVictims 关闭超时的靶机
func closeTimeoutVictims(c *cron.Cron) {
	function := exec("CloseTimeoutVictims", func() {
		repo := db.InitVictimRepo(db.DB)
		victims, _, ok, _ := repo.List(-1, -1, db.GetOptions{
			Preloads: map[string]db.GetOptions{
				"Team": {
					Selects: []string{"id", "contest_id"},
					Preloads: map[string]db.GetOptions{
						"Contest": {Selects: []string{"id", "start", "duration"}},
					},
				},
			},
		})
		if !ok {
			return
		}
		for _, victim := range victims {
			if victim.Start.Add(victim.Duration).Before(time.Now()) || (victim.TeamID.Valid && victim.Team.Contest.IsOver()) {
				service.StopVictim(db.DB, victim)
			}
		}
	})
	function()
	c.Schedule(cron.Every(5*time.Minute), cron.FuncJob(function))
}

// closeUnCtrlVictims 关闭数据库中记录关闭, 但仍在运行的靶机
func closeUnCtrlVictims(c *cron.Cron) {
	function := exec("CloseUnCtrlVictims", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		pods, ok, msg := k8s.GetPodList(ctx)
		cancel()
		if !ok {
			log.Logger.Warningf("Failed to get Victim %s", msg)
			return
		}
		idL := make([]string, 0)
		victimRepo := db.InitVictimRepo(db.DB)
		for _, pod := range pods.Items {
			for key := range pod.Labels {
				if key == "victim_id" {
					if slices.Contains(idL, pod.Labels[key]) {
						continue
					}
					victimID, err := strconv.Atoi(pod.Labels[key])
					if err != nil {
						continue
					}
					_, ok, _ = victimRepo.GetByID(uint(victimID), db.GetOptions{Selects: []string{"id"}})
					if !ok {
						idL = append(idL, pod.Labels[key])
					}
				}
			}
		}
		for _, id := range idL {
			ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
			k8s.DeletePodList(ctx, map[string]string{"victim_id": id})
			cancel()
		}
	})
	function()
	c.Schedule(cron.Every(time.Hour), cron.FuncJob(function))
}
