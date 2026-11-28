package cron

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"

	corev1 "k8s.io/api/core/v1"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/service"
	"CBCTF/internal/task"
)

// closeTimeoutVictimsTask 关闭超时的靶机
func closeTimeoutVictimsTask() model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	root := db.CronDB.WithContext(ctx)
	repo := db.InitVictimRepo(root)
	victims, _, ret := repo.List(-1, -1, db.GetOptions{Conditions: map[string]any{"status": model.RunningVictimStatus}})
	if !ret.OK {
		return ret
	}
	batch := model.NewBatch(len(victims))
	for _, victim := range victims {
		if ctx.Err() != nil {
			return batch.Result(ctx)
		}
		key := fmt.Sprint(victim.ID)
		if victimExpired(victim, time.Now()) {
			if ret = service.ForceStopVictim(root, victim); ret.OK {
				prometheus.RecordVictimTransition("expire", true)
				batch.Success(key, "queued")
				log.Logger.Infof(
					"Timeout victim stop queued: victim_id=%d team_id=%d challenge_id=%d expired_at=%s",
					victim.ID, victim.TeamID.V, victim.ChallengeID, victim.Start.Add(victim.Duration).Format(time.RFC3339),
				)
			} else {
				prometheus.RecordVictimTransition("expire", false)
				batch.Fail(key, "stop", ret)
				log.Logger.Warningf("Failed to queue timeout victim stop: victim_id=%d reason=%v", victim.ID, ret)
				if model.BatchDependencyFailed(ret) {
					return batch.Result(ctx)
				}
			}
		} else {
			batch.Skip(key, "not_expired")
		}
	}
	if ret = requeueStuckVictims(root, batch); !ret.OK {
		return ret
	}
	return batch.Result(ctx)
}

// stuckVictimGrace 是实例卡在 terminating 之后重新投递停止任务的等待时间。
// 停止任务重试耗尽实例就会一直卡住,资源留在集群里,这里兜底重投。
const stuckVictimGrace = 5 * time.Minute

// requeueStuckVictims 重新投递卡在 terminating 的实例的停止任务。
func requeueStuckVictims(root *gorm.DB, batch *model.BatchResult) model.RetVal {
	victims, _, ret := db.InitVictimRepo(root).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"status": model.TerminatingVictimStatus},
	})
	if !ret.OK {
		return ret
	}
	for _, victim := range victims {
		if time.Since(victim.UpdatedAt) < stuckVictimGrace {
			continue
		}
		key := fmt.Sprint(victim.ID)
		if err := task.EnqueueStopVictimTask(victim); err != nil {
			batch.Fail(key, "requeue", model.FailRetVal(i18n.Task.EnqueueError, "Victim", key, err))
			continue
		}
		batch.Success(key, "requeued")
		log.Logger.Infof("Stuck victim stop requeued: victim_id=%d team_id=%d challenge_id=%d", victim.ID, victim.TeamID.V, victim.ChallengeID)
	}
	return model.SuccessRetVal()
}

// closeUnCtrlVictimsTask 关闭数据库中记录关闭, 但仍在运行的靶机
func closeUnCtrlVictimsTask() model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	listCtx, listCancel := context.WithTimeout(ctx, time.Minute)
	roots, err := k8s.ListVictimRoots(listCtx)
	listCancel()
	if err != nil {
		return model.FailRetVal(i18n.K8S.GetError, "VictimRoot", "", err)
	}
	var orphans []model.Victim
	for _, victim := range roots {
		_, ret := db.InitVictimRepo(db.CronDB.WithContext(ctx)).GetByID(victim.ID)
		if !ret.OK && ret.Msg != i18n.Model.NotFound {
			return ret
		}
		if !ret.OK {
			orphans = append(orphans, victim)
		}
	}
	batch := model.NewBatch(len(orphans))
	for _, victim := range orphans {
		if ctx.Err() != nil {
			return batch.Result(ctx)
		}
		itemCtx, itemCancel := context.WithTimeout(ctx, 2*time.Minute)
		ret := k8s.StopVictim(itemCtx, victim)
		itemCancel()
		if !ret.OK {
			batch.Fail(fmt.Sprint(victim.ID), "cleanup", ret)
			continue
		}
		batch.Success(fmt.Sprint(victim.ID), "deleted")
		log.Logger.Infof("Deleted uncontrolled victim resource tree: victim_id=%d", victim.ID)
	}
	return batch.Result(ctx)
}

func victimExpired(victim model.Victim, now time.Time) bool {
	if victim.Status != model.RunningVictimStatus {
		return false
	}
	// 靶场实例按租约到期回收,池中未租出的实例不计时。
	if victim.IsRange() {
		return victim.TeamID.Valid && !victim.LeaseExpireAt.IsZero() && victim.LeaseExpireAt.Before(now)
	}
	return !victim.Start.IsZero() && victim.Start.Add(victim.Duration).Before(now)
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
