package task

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	startVictimTaskType = "tasks:victim:start"
	stopVictimTaskType  = "tasks:victim:stop"
)

type StartVictimPayload struct {
	Victim model.Victim
}

func EnqueueStartVictimTask(victim model.Victim) (*asynq.TaskInfo, error) {
	payload, err := msgpack.Marshal(StartVictimPayload{victim})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(startVictimTaskType, payload)
	return enqueueTask(startVictimTaskType, task, asynq.MaxRetry(0), asynq.Timeout(4*time.Minute))
}

func HandleStartVictimTask(ctx context.Context, t *asynq.Task) error {
	var payload StartVictimPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return db.WithWorkloadLock(ctx, db.WorkloadLockDB, "victim", payload.Victim.ID, func() error {
		victim := payload.Victim
		log.Logger.Debugf("Start victim task received: victim_id=%d user_id=%d team_id=%d challenge_id=%d pods=%d", victim.ID, victim.UserID, victim.TeamID.V, victim.ChallengeID, len(victim.Pods))
		cleanupQueued := false
		claimed := false
		cleanupFailedStart := func(reason error) error {
			victimRepo := db.InitVictimRepo(db.TaskDB)
			// The stop worker shares our lock; persist partial allocation even if a
			// stop request already changed status while Kubernetes was provisioning.
			_ = victimRepo.Update(victim.ID, db.UpdateVictimOptions{Resources: &victim.Resources, ExposedEndpoints: &victim.ExposedEndpoints})
			expectedStatus := victim.Status
			victim.Status = model.TerminatingVictimStatus
			if ret := victimRepo.UpdateIfStatus(victim.ID, expectedStatus, db.UpdateVictimOptions{Status: new(model.TerminatingVictimStatus), Resources: &victim.Resources, ExposedEndpoints: &victim.ExposedEndpoints}); !ret.OK {
				log.Logger.Warningf("Failed to mark victim terminating after start failure: victim_id=%d reason=%s", victim.ID, ret.Msg)
			}
			if enqueueErr := EnqueueStopVictimTask(victim); enqueueErr == nil {
				cleanupQueued = true
				return reason
			} else {
				log.Logger.Warningf("Failed to enqueue victim cleanup after start failure: victim_id=%d error=%v", victim.ID, enqueueErr)
			}

			cleanupCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			ret := k8s.StopVictim(cleanupCtx, victim)
			cancel()
			if !ret.OK {
				return fmt.Errorf("%w; cleanup enqueue failed and synchronous cleanup failed: %s", reason, ret.Msg)
			}
			if ret = victimRepo.Delete(victim.ID); !ret.OK {
				return fmt.Errorf("%w; synchronous cleanup deleted k8s resources but failed to delete victim: %s", reason, ret.Msg)
			}
			cleanupQueued = true
			return reason
		}
		err := func() error {
			podRepo := db.InitPodRepo(db.TaskDB)
			victimRepo := db.InitVictimRepo(db.TaskDB)
			currentVictim, ret := victimRepo.GetByID(victim.ID, db.GetOptions{
				Preloads: map[string]db.GetOptions{"Pods": {}},
			})
			if !ret.OK {
				if ret.Msg == i18n.Model.NotFound {
					log.Logger.Debugf("Start victim skipped: victim_id=%d no longer exists", victim.ID)
					return nil
				}
				return fmt.Errorf("get victim failed: %s", ret.Msg)
			}
			if currentVictim.Status != model.WaitingVictimStatus {
				log.Logger.Infof("Start victim skipped: victim_id=%d status=%s", victim.ID, currentVictim.Status)
				return nil
			}
			victim = currentVictim
			victim.Resources.ReadyDeadline = time.Now().Add(4 * time.Minute)
			if ret = victimRepo.UpdateIfStatus(victim.ID, model.WaitingVictimStatus, db.UpdateVictimOptions{Status: new(model.PendingVictimStatus), Resources: &victim.Resources}); !ret.OK {
				if ret.Msg == i18n.Model.Victim.NotStartable {
					log.Logger.Infof("Start victim skipped: victim_id=%d status changed before provisioning", victim.ID)
					return nil
				}
				return fmt.Errorf("update victim failed: %s", ret.Msg)
			}
			claimed = true
			victim.Status = model.PendingVictimStatus
			basePodCount := len(victim.Pods)
			log.Logger.Infof("Starting victim provisioning: victim_id=%d user_id=%d team_id=%d challenge_id=%d", victim.ID, victim.UserID, victim.TeamID.V, victim.ChallengeID)
			victim, ret = k8s.StartVictim(ctx, victim)
			if !ret.OK {
				return taskResourceError("start victim failed", ret)
			}
			if len(victim.Pods) > basePodCount {
				frpcPods := victim.Pods[basePodCount:]
				frpcPodRecords := make([]model.Pod, 0, len(frpcPods))
				for _, frpcPod := range frpcPods {
					frpcPodRecords = append(frpcPodRecords, model.Pod{
						VictimID: victim.ID,
						Name:     frpcPod.Name,
					})
				}
				persistedFrpcPods, ret := podRepo.CreateBatch(frpcPodRecords)
				if !ret.OK {
					return fmt.Errorf("create frpc pod failed: %s", ret.Msg)
				}
				victim.Pods = append(append([]model.Pod(nil), victim.Pods[:basePodCount]...), persistedFrpcPods...)
			}
			victim.Resources.Submitted = true
			ret = victimRepo.UpdateIfStatus(victim.ID, model.PendingVictimStatus, db.UpdateVictimOptions{
				Spec:             &victim.Spec,
				Resources:        &victim.Resources,
				Endpoints:        &victim.Endpoints,
				ExposedEndpoints: &victim.ExposedEndpoints,
			})
			if !ret.OK {
				return fmt.Errorf("update victim after start failed: %s", ret.Msg)
			}
			log.Logger.Infof(
				"Victim resources submitted: victim_id=%d user_id=%d team_id=%d challenge_id=%d endpoints=%d exposed_endpoints=%d",
				victim.ID, victim.UserID, victim.TeamID.V, victim.ChallengeID, len(victim.Endpoints), len(victim.ExposedEndpoints),
			)
			return nil
		}()
		if err != nil && claimed && !cleanupQueued {
			err = cleanupFailedStart(err)
		}
		return err
	})
}

type StopVictimPayload struct {
	Victim model.Victim
}

func EnqueueStopVictimTask(victim model.Victim) error {
	payload, err := msgpack.Marshal(StopVictimPayload{victim})
	if err != nil {
		return err
	}
	task := asynq.NewTask(stopVictimTaskType, payload)
	_, err = enqueueTask(stopVictimTaskType, task, asynq.MaxRetry(3), asynq.Timeout(2*time.Minute))
	return err
}

func HandleStopVictimTask(ctx context.Context, t *asynq.Task) error {
	var payload StopVictimPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return db.WithWorkloadLock(ctx, db.WorkloadLockDB, "victim", payload.Victim.ID, func() error {
		victimRepo := db.InitVictimRepo(db.TaskDB)
		victim, ret := victimRepo.GetByID(payload.Victim.ID, db.GetOptions{Preloads: map[string]db.GetOptions{"Pods": {}}})
		if !ret.OK {
			if ret.Msg == i18n.Model.NotFound {
				return nil
			}
			return fmt.Errorf("get victim failed: %s", ret.Msg)
		}
		if len(payload.Victim.ExposedEndpoints) > 0 && len(victim.ExposedEndpoints) == 0 {
			victim.ExposedEndpoints = payload.Victim.ExposedEndpoints
		}
		if len(victim.Pods) == 0 {
			victim.Pods = payload.Victim.Pods
		}
		log.Logger.Infof("Stopping victim: victim_id=%d user_id=%d team_id=%d challenge_id=%d", victim.ID, victim.UserID, victim.TeamID.V, victim.ChallengeID)
		ret = k8s.StopVictim(ctx, victim)
		if !ret.OK {
			return fmt.Errorf("stop victim failed: %s", ret.Msg)
		}
		if _, err := EnqueueLoadTrafficTask(victim); err != nil {
			log.Logger.Warningf("Failed to enqueue load traffic task: victim_id=%d user_id=%d team_id=%d challenge_id=%d error=%v", victim.ID, victim.UserID, victim.TeamID.V, victim.ChallengeID, err)
		}
		ret = db.WithTransactionDB(db.TaskDB, func(tx *db.Tx) model.RetVal {
			if ret = db.InitVictimRepo(tx).Update(victim.ID, db.UpdateVictimOptions{
				Duration: new(time.Now().Sub(victim.Start)),
			}); !ret.OK {
				return model.RetVal{Msg: fmt.Sprintf("update victim failed: %s", ret.Msg)}
			}
			if ret = db.InitVictimRepo(tx).Delete(victim.ID); !ret.OK {
				return model.RetVal{Msg: fmt.Sprintf("delete victim failed: %s", ret.Msg)}
			}
			return model.SuccessRetVal()
		})
		if !ret.OK {
			return fmt.Errorf("%s", ret.Msg)
		}
		log.Logger.Infof("Victim stopped: victim_id=%d user_id=%d team_id=%d challenge_id=%d", victim.ID, victim.UserID, victim.TeamID.V, victim.ChallengeID)
		return nil
	})
}
