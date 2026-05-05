package task

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
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
	info, err := client.Enqueue(task, asynq.Queue(startVictimTaskType), asynq.MaxRetry(0), asynq.Timeout(4*time.Minute))
	if err == nil {
		prometheus.RecordTaskEnqueued(startVictimTaskType)
	}
	return info, err
}

func HandleStartVictimTask(_ context.Context, t *asynq.Task) error {
	var payload StartVictimPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	victim := payload.Victim
	log.Logger.Debugf("Start victim task received: victim_id=%d user_id=%d team_id=%d challenge_id=%d pods=%d", victim.ID, victim.UserID, victim.TeamID.V, victim.ChallengeID, len(victim.Pods))
	cleanupQueued := false
	err := func() error {
		podRepo := db.InitPodRepo(db.DB)
		victimRepo := db.InitVictimRepo(db.DB)
		currentVictim, ret := victimRepo.GetByID(victim.ID)
		if !ret.OK {
			if ret.Msg == i18n.Model.NotFound {
				log.Logger.Debugf("Start victim skipped: victim_id=%d no longer exists", victim.ID)
				return nil
			}
			return fmt.Errorf("get victim failed: %s", ret.Msg)
		}
		if currentVictim.Status == model.TerminatingVictimStatus {
			log.Logger.Infof("Start victim skipped: victim_id=%d is terminating", victim.ID)
			return nil
		}
		if ret = victimRepo.Update(victim.ID, db.UpdateVictimOptions{Status: new(model.PendingVictimStatus)}); !ret.OK {
			return fmt.Errorf("update victim failed: %s", ret.Msg)
		}
		log.Logger.Infof("Starting victim provisioning: victim_id=%d user_id=%d team_id=%d challenge_id=%d", victim.ID, victim.UserID, victim.TeamID.V, victim.ChallengeID)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()
		victim, ret = k8s.StartVictim(ctx, victim)
		if !ret.OK {
			return fmt.Errorf("start victim failed: %s", ret.Msg)
		}
		basePodCount := len(payload.Victim.Pods)
		if len(victim.Pods) > basePodCount {
			frpcPods := victim.Pods[basePodCount:]
			persistedFrpcPods := make([]model.Pod, 0, len(frpcPods))
			for _, frpcPod := range frpcPods {
				p, ret := podRepo.Create(db.CreatePodOptions{
					VictimID: victim.ID,
					Name:     frpcPod.Name,
				})
				if !ret.OK {
					return fmt.Errorf("create frpc pod failed: %s", ret.Msg)
				}
				persistedFrpcPods = append(persistedFrpcPods, p)
			}
			victim.Pods = append(append([]model.Pod(nil), victim.Pods[:basePodCount]...), persistedFrpcPods...)
		}
		ret = victimRepo.Update(victim.ID, db.UpdateVictimOptions{
			Spec:             &victim.Spec,
			Resources:        &victim.Resources,
			Endpoints:        &victim.Endpoints,
			ExposedEndpoints: &victim.ExposedEndpoints,
			Start:            new(time.Now()),
			Status:           new(model.RunningVictimStatus),
		})
		if !ret.OK {
			_, err := EnqueueStopVictimTask(victim)
			cleanupQueued = err == nil
			return err
		}
		log.Logger.Infof(
			"Victim is running: victim_id=%d user_id=%d team_id=%d challenge_id=%d endpoints=%d exposed_endpoints=%d",
			victim.ID, victim.UserID, victim.TeamID.V, victim.ChallengeID, len(victim.Endpoints), len(victim.ExposedEndpoints),
		)
		return nil
	}()
	if err != nil && !cleanupQueued {
		if _, enqueueErr := EnqueueStopVictimTask(victim); enqueueErr != nil {
			log.Logger.Warningf("Failed to enqueue victim cleanup after start failure: victim_id=%d error=%v", victim.ID, enqueueErr)
		}
	}
	return err
}

type StopVictimPayload struct {
	Victim model.Victim
}

func EnqueueStopVictimTask(victim model.Victim) (*asynq.TaskInfo, error) {
	payload, err := msgpack.Marshal(StopVictimPayload{victim})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(stopVictimTaskType, payload)
	info, err := client.Enqueue(task, asynq.Queue(stopVictimTaskType), asynq.MaxRetry(3), asynq.Timeout(2*time.Minute))
	if err == nil {
		prometheus.RecordTaskEnqueued(stopVictimTaskType)
	}
	return info, err
}

func HandleStopVictimTask(ctx context.Context, t *asynq.Task) error {
	var payload StopVictimPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	victimRepo := db.InitVictimRepo(db.DB)
	victim, ret := victimRepo.GetByID(payload.Victim.ID)
	if !ret.OK {
		if ret.Msg == i18n.Model.NotFound {
			log.Logger.Debugf("Stop victim skipped: victim_id=%d no longer exists", payload.Victim.ID)
			return nil
		}
		return fmt.Errorf("get victim failed: %s", ret.Msg)
	}
	log.Logger.Infof("Stopping victim: victim_id=%d user_id=%d team_id=%d challenge_id=%d", victim.ID, victim.UserID, victim.TeamID.V, victim.ChallengeID)
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	ret = k8s.StopVictim(ctx, victim)
	if !ret.OK {
		return fmt.Errorf("stop victim failed: %s", ret.Msg)
	}
	ret = db.WithTransaction(func(tx *db.Tx) model.RetVal {
		if ret := db.InitVictimRepo(tx).Update(victim.ID, db.UpdateVictimOptions{
			Duration: new(time.Now().Sub(victim.Start)),
		}); !ret.OK {
			return model.RetVal{Msg: fmt.Sprintf("update victim failed: %s", ret.Msg)}
		}
		if ret := db.InitVictimRepo(tx).Delete(victim.ID); !ret.OK {
			return model.RetVal{Msg: fmt.Sprintf("delete victim failed: %s", ret.Msg)}
		}
		return model.SuccessRetVal()
	})
	if !ret.OK {
		return fmt.Errorf("%s", ret.Msg)
	}
	log.Logger.Infof("Victim stopped: victim_id=%d user_id=%d team_id=%d challenge_id=%d", victim.ID, victim.UserID, victim.TeamID.V, victim.ChallengeID)
	return nil
}
