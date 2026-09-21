package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"context"
	corev1 "k8s.io/api/core/v1"
	"strconv"
	"time"
)

// closeTimeoutVictimsTask 关闭超时的靶机
func closeTimeoutVictimsTask() model.RetVal {
	repo := db.InitVictimRepo(db.CronDB)
	victims, _, ret := repo.List(-1, -1, db.GetOptions{Conditions: map[string]any{"status": model.RunningVictimStatus}})
	if !ret.OK {
		return ret
	}
	for _, victim := range victims {
		if victimExpired(victim, time.Now()) {
			if ret = service.ForceStopVictim(db.CronDB, victim); ret.OK {
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
	pods, ret := k8s.ListPods(ctx)
	cancel()
	if !ret.OK {
		return ret
	}
	if pods == nil {
		return model.SuccessRetVal()
	}
	idL, ret := orphanVictimIDs(pods.Items, func(id uint) model.RetVal {
		_, ret := db.InitVictimRepo(db.CronDB).GetByID(id)
		return ret
	})
	if !ret.OK {
		return ret
	}
	for _, id := range idL {
		ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
		if ret = k8s.DeletePodCollection(ctx, map[string]string{"victim_id": id}); ret.OK {
			log.Logger.Infof("Deleted uncontrolled victim pods: victim_id=%s", id)
		} else {
			log.Logger.Warningf("Failed to delete uncontrolled victim pods: victim_id=%s reason=%s", id, ret.Msg)
		}
		cancel()
	}
	return model.SuccessRetVal()
}

func victimExpired(victim model.Victim, now time.Time) bool {
	return victim.Status == model.RunningVictimStatus && !victim.Start.IsZero() && victim.Start.Add(victim.Duration).Before(now)
}

// Build the complete deletion candidate set before any destructive API call.
// A database failure invalidates the entire scan, not just the current Pod.
func orphanVictimIDs(pods []corev1.Pod, lookup func(uint) model.RetVal) ([]string, model.RetVal) {
	ids := make([]string, 0)
	seen := make(map[string]bool)
	for _, pod := range pods {
		label := pod.Labels["victim_id"]
		if seen[label] {
			continue
		}
		seen[label] = true
		id, err := strconv.ParseUint(label, 10, 64)
		if err != nil || id == 0 || uint64(uint(id)) != id {
			continue
		}
		ret := lookup(uint(id))
		if !ret.OK && ret.Msg != i18n.Model.NotFound {
			return nil, ret
		}
		if !ret.OK {
			ids = append(ids, label)
		}
	}
	return ids, model.SuccessRetVal()
}
