package task

import (
	"context"
	"fmt"
	"time"

	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

var readinessCancel context.CancelFunc
var readinessDone chan struct{}

func startReadinessController() {
	ctx, cancel := context.WithCancel(context.Background())
	readinessCancel, readinessDone = cancel, make(chan struct{})
	go func() {
		defer close(readinessDone)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
			victims, _, ret := db.InitVictimRepo(db.TaskDB.WithContext(ctx)).List(-1, -1, db.GetOptions{
				Conditions: map[string]any{"status": model.PendingVictimStatus},
			})
			if !ret.OK {
				log.Logger.Warningf("List pending victims failed: %s", ret.Msg)
				continue
			}
			for _, victim := range victims {
				checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
				err := db.TryWithWorkloadLock(checkCtx, db.WorkloadLockDB, "victim", victim.ID, func() error {
					return reconcileVictim(checkCtx, victim.ID)
				})
				if err != nil && checkCtx.Err() == nil {
					log.Logger.Warningf("Reconcile victim %d: %v", victim.ID, err)
				}
				cancel()
			}
		}
	}()
}

func reconcileVictim(ctx context.Context, id uint) error {
	repo := db.InitVictimRepo(db.TaskDB.WithContext(ctx))
	victim, ret := repo.GetByID(id, db.GetOptions{Preloads: map[string]db.GetOptions{"Pods": {}}})
	if !ret.OK || victim.Status != model.PendingVictimStatus {
		return nil
	}
	ready, err := k8s.VictimReady(ctx, &victim)
	expired := time.Now().After(victim.Resources.ReadyDeadline)
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err != nil || expired {
		log.Logger.Warningf("Victim readiness failed: victim_id=%d expired=%t error=%v", id, expired, err)
		// Queue under the same workload lock before changing status. A crash in
		// between leaves either a retryable pending record or durable cleanup.
		if err := EnqueueStopVictimTask(victim); err != nil {
			return err
		}
		ret = repo.UpdateIfStatus(id, model.PendingVictimStatus, db.UpdateVictimOptions{Status: new(model.TerminatingVictimStatus)})
		if !ret.OK {
			return fmt.Errorf("mark terminating: %s", ret.Msg)
		}
		return nil
	}
	if !ready {
		return nil
	}
	ret = repo.UpdateIfStatus(id, model.PendingVictimStatus, db.UpdateVictimOptions{
		Status:           new(model.RunningVictimStatus),
		Start:            new(time.Now()),
		Endpoints:        &victim.Endpoints,
		ExposedEndpoints: &victim.ExposedEndpoints,
	})
	if !ret.OK {
		return fmt.Errorf("mark running: %s", ret.Msg)
	}
	log.Logger.Infof("Victim is running: victim_id=%d", id)
	return nil
}
