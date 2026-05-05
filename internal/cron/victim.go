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
			if ret = service.ForceStopVictim(db.DB, victim); ret.OK {
				log.Logger.Infof(
					"Timeout victim stop queued: victim_id=%d team_id=%d challenge_id=%d expired_at=%s",
					victim.ID, victim.TeamID.V, victim.ChallengeID, victim.Start.Add(victim.Duration).Format(time.RFC3339),
				)
			} else {
				log.Logger.Warningf("Failed to queue timeout victim stop: victim_id=%d reason=%s", victim.ID, ret.Msg)
			}
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
		if ret = k8s.DeletePodList(ctx, map[string]string{"victim_id": id}); ret.OK {
			log.Logger.Infof("Deleted uncontrolled victim pods: victim_id=%s", id)
		} else {
			log.Logger.Warningf("Failed to delete uncontrolled victim pods: victim_id=%s reason=%s", id, ret.Msg)
		}
		cancel()
	}
	return model.SuccessRetVal()
}
