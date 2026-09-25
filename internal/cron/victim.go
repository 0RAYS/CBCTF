package cron

import (
	"context"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
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
	roots, err := k8s.ListVictimRoots(ctx)
	cancel()
	if err != nil {
		return model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "VictimRoot", "Error": err.Error()}}
	}
	var orphans []model.Victim
	for _, victim := range roots {
		_, ret := db.InitVictimRepo(db.CronDB).GetByID(victim.ID)
		if !ret.OK && ret.Msg != i18n.Model.NotFound {
			return ret
		}
		if !ret.OK {
			orphans = append(orphans, victim)
		}
	}
	for _, victim := range orphans {
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute)
		ret := k8s.StopVictim(ctx, victim)
		cancel()
		if !ret.OK {
			return ret
		}
		log.Logger.Infof("Deleted uncontrolled victim resource tree: victim_id=%d", victim.ID)
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
