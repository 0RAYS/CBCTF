package task

import (
	"CBCTF/internal/config"
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
	info, err := client.Enqueue(task, asynq.Queue(startVictimTaskType), asynq.MaxRetry(0), asynq.Timeout(2*time.Minute))
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
	podRepo := db.InitPodRepo(db.DB)
	victimRepo := db.InitVictimRepo(db.DB)
	if _, ret := victimRepo.GetByID(victim.ID); !ret.OK {
		if ret.Msg == i18n.Model.NotFound {
			log.Logger.Infof("The Victim %d may have already been stopped", victim.ID)
			return nil
		}
		return fmt.Errorf("get victim failed: %s", ret.Msg)
	}
	if ret := victimRepo.Update(victim.ID, db.UpdateVictimOptions{Status: new(model.PendingVictimStatus)}); !ret.OK {
		return fmt.Errorf("update victim failed: %s", ret.Msg)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	ipExposesMap, ret := k8s.StartVictim(ctx, victim)
	if !ret.OK {
		return fmt.Errorf("start victim failed: %s", ret.Msg)
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
			return fmt.Errorf("create frpc failed: %s", ret.Msg)
		}
		for _, frpcPodName := range frpc {
			p, ret := podRepo.Create(db.CreatePodOptions{
				VictimID: victim.ID,
				Name:     frpcPodName,
			})
			if !ret.OK {
				return fmt.Errorf("create frpc pod failed: %s", ret.Msg)
			}
			victim.Pods = append(victim.Pods, p)
		}
	}
	ret = victimRepo.Update(victim.ID, db.UpdateVictimOptions{
		VPC:              &victim.VPC,
		Endpoints:        &victim.Endpoints,
		ExposedEndpoints: &victim.ExposedEndpoints,
		Start:            new(time.Now()),
		Status:           new(model.RunningVictimStatus),
	})
	if !ret.OK {
		victim, ret = db.InitVictimRepo(db.DB).HasAliveVictim(victim.TeamID.V, victim.ChallengeID)
		if !ret.OK {
			return fmt.Errorf("check victim failed: %s", ret.Msg)
		}
		_, err := EnqueueStopVictimTask(victim)
		return err
	}
	return nil
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
	victim := payload.Victim
	if victim.Status == model.PendingVictimStatus {
		log.Logger.Infof("The Victim %d is pending, skip it...", victim.ID)
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	ret := k8s.StopVictim(ctx, victim)
	if !ret.OK {
		return fmt.Errorf("stop victim failed: %s", ret.Msg)
	}
	tx := db.DB.Begin()
	if ret = db.InitVictimRepo(tx).Update(victim.ID, db.UpdateVictimOptions{
		Duration: new(time.Now().Sub(victim.Start)),
	}); !ret.OK {
		tx.Rollback()
		return fmt.Errorf("update victim failed: %s", ret.Msg)
	}
	if ret = db.InitVictimRepo(tx).Delete(victim.ID); !ret.OK {
		tx.Rollback()
		return fmt.Errorf("delete victim failed: %s", ret.Msg)
	}
	tx.Commit()
	return nil
}
