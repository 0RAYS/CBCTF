package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"context"
	"slices"
	"strconv"
	"time"
)

// closeTimeoutVictimsTask 关闭超时的靶机
func closeTimeoutVictimsTask() model.RetVal {
	repo := db.InitVictimRepo(db.DB)
	victims, _, ret := repo.List(-1, -1)
	if !ret.OK {
		return ret
	}
	for _, victim := range victims {
		if victim.Start.Add(victim.Duration).Before(time.Now()) {
			service.StopVictim(db.DB, victim)
		}
	}
	return model.SuccessRetVal()
}

// closeUnCtrlVictimsTask 关闭数据库中记录关闭, 但仍在运行的靶机
func closeUnCtrlVictimsTask() model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	pods, ret := k8s.GetPodList(ctx)
	cancel()
	if !ret.OK {
		log.Logger.Warningf("Failed to get Victim %v", ret)
		return ret
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
				_, ret = victimRepo.GetByID(uint(victimID))
				if !ret.OK {
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
	return model.SuccessRetVal()
}
