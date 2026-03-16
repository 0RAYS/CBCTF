package task

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	StartVictimTaskType = "tasks:victim:start"
	StopVictimTaskType  = "tasks:victim:stop"
)

type StartVictimPayload struct {
	Victim model.Victim
}

func EnqueueStartVictimTask(victim model.Victim) (*asynq.TaskInfo, error) {
	payload, err := msgpack.Marshal(StartVictimPayload{victim})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(StartVictimTaskType, payload)
	info, err := client.Enqueue(task, asynq.Queue(queueForTask(StartVictimTaskType)), asynq.MaxRetry(0), asynq.Timeout(2*time.Minute))
	if err == nil {
		prometheus.RecordTaskEnqueued(StartVictimTaskType)
	}
	return info, err
}

func HandleStartVictimTask(ctx context.Context, t *asynq.Task) error {
	var payload StartVictimPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	victim := payload.Victim
	podRepo := db.InitPodRepo(db.DB)
	victimRepo := db.InitVictimRepo(db.DB)
	ret := func() model.RetVal {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()
		ipExposesMap, ret := k8s.StartVictim(ctx, victim)
		if !ret.OK {
			return ret
		}
		for ip, exposes := range ipExposesMap {
			for _, expose := range exposes {
				victim.Endpoints = append(victim.Endpoints, model.Endpoint{
					IP:       ip,
					Port:     expose.Port,
					Protocol: expose.Protocol,
				})
			}
		}
		victim.ExposedEndpoints = victim.Endpoints
		if config.Env.K8S.Frp.On {
			var frpc []string
			victim.ExposedEndpoints, frpc, ret = k8s.CreateFrpc(ctx, victim)
			if !ret.OK {
				return ret
			}
			for _, frpcPodName := range frpc {
				p, ret := podRepo.Create(db.CreatePodOptions{
					VictimID: victim.ID,
					Name:     frpcPodName,
				})
				if !ret.OK {
					return ret
				}
				victim.Pods = append(victim.Pods, p)
			}
		}
		if ret = victimRepo.Update(victim.ID, db.UpdateVictimOptions{
			VPC:              &victim.VPC,
			Endpoints:        &victim.Endpoints,
			ExposedEndpoints: &victim.ExposedEndpoints,
			Start:            new(time.Now()),
			Status:           new(model.RunningVictimStatus),
		}); !ret.OK {
			return ret
		}
		return model.SuccessRetVal()
	}()
	err := error(nil)
	if !ret.OK {
		victim, ret = db.InitVictimRepo(db.DB).HasAliveVictim(victim.TeamID.V, victim.ChallengeID)
		if !ret.OK {
			return fmt.Errorf("start victim fail: %s", ret.Msg)
		}
		_, err = EnqueueStopVictimTask(victim)
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
	task := asynq.NewTask(StopVictimTaskType, payload)
	info, err := client.Enqueue(task, asynq.Queue(queueForTask(StopVictimTaskType)), asynq.MaxRetry(0), asynq.Timeout(2*time.Minute))
	if err == nil {
		prometheus.RecordTaskEnqueued(StopVictimTaskType)
	}
	return info, err
}

func HandleStopVictimTask(ctx context.Context, t *asynq.Task) error {
	var payload StopVictimPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	victim := payload.Victim
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	ret := k8s.StopVictim(ctx, victim)
	if !ret.OK {
		return fmt.Errorf("stop victim fail: %s", ret.Msg)
	}
	tx := db.DB.Begin()
	if ret = db.InitVictimRepo(tx).Update(victim.ID, db.UpdateVictimOptions{
		Duration: new(time.Now().Sub(victim.Start)),
	}); !ret.OK {
		tx.Rollback()
		return fmt.Errorf("stop victim fail, update victim fail %s", ret.Msg)
	}
	if ret = db.InitVictimRepo(tx).Delete(victim.ID); !ret.OK {
		tx.Rollback()
		return fmt.Errorf("stop victim fail, delete victim fail %s", ret.Msg)
	}
	tx.Commit()
	return nil
}
