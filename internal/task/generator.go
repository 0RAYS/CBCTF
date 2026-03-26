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
	startGeneratorTaskType = "tasks:generator:start"
	stopGeneratorTaskType  = "tasks:generator:stop"
)

type StartGeneratorPayload struct {
	Challenge model.Challenge
	Generator model.Generator
}

func EnqueueStartGeneratorTask(challenge model.Challenge, generator model.Generator) (*asynq.TaskInfo, error) {
	payload, err := msgpack.Marshal(StartGeneratorPayload{challenge, generator})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(startGeneratorTaskType, payload)
	info, err := client.Enqueue(task, asynq.Queue(startGeneratorTaskType), asynq.MaxRetry(0), asynq.Timeout(2*time.Minute))
	if err != nil {
		prometheus.RecordTaskEnqueued(startGeneratorTaskType)
	}
	return info, err
}

func HandleStartGeneratorTask(_ context.Context, t *asynq.Task) error {
	var payload StartGeneratorPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	err := func() error {
		challenge := payload.Challenge
		generator := payload.Generator
		generatorRepo := db.InitGeneratorRepo(db.DB)
		currentGenerator, ret := generatorRepo.GetByID(generator.ID)
		if !ret.OK {
			if ret.Msg == i18n.Model.NotFound {
				log.Logger.Infof("The Generator %d may have already been stopped", generator.ID)
				return nil
			}
			return fmt.Errorf("get generator failed: %s", ret.Msg)
		}
		if currentGenerator.Status == model.TerminatingGeneratorStatus {
			log.Logger.Infof("The Generator %d is terminating, skip start...", generator.ID)
			return nil
		}
		if ret := generatorRepo.Update(generator.ID, db.UpdateGeneratorOptions{Status: new(model.PendingGeneratorStatus)}); !ret.OK {
			return fmt.Errorf("update generator failed: %s", ret.Msg)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		_, ret = k8s.StartGenerator(ctx, challenge, generator)
		cancel()
		if !ret.OK {
			_, err := EnqueueStopGeneratorTask(generator)
			return err
		}
		ret = generatorRepo.Update(generator.ID, db.UpdateGeneratorOptions{Status: new(model.RunningGeneratorStatus)})
		if !ret.OK {
			return fmt.Errorf("update generator failed: %s", ret.Msg)
		}
		return nil
	}()
	if err != nil {
		_, _ = EnqueueStopGeneratorTask(payload.Generator)
	}
	return err
}

type StopGeneratorPayload struct {
	Generator model.Generator
}

func EnqueueStopGeneratorTask(generator model.Generator) (*asynq.TaskInfo, error) {
	payload, err := msgpack.Marshal(StopGeneratorPayload{generator})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(stopGeneratorTaskType, payload)
	info, err := client.Enqueue(task, asynq.Queue(stopGeneratorTaskType), asynq.MaxRetry(3), asynq.Timeout(2*time.Minute))
	if err != nil {
		prometheus.RecordTaskEnqueued(stopGeneratorTaskType)
	}
	return info, err
}

func HandleStopGeneratorTask(ctx context.Context, t *asynq.Task) error {
	var payload StopGeneratorPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	generatorRepo := db.InitGeneratorRepo(db.DB)
	generator, ret := generatorRepo.GetByID(payload.Generator.ID)
	if !ret.OK {
		if ret.Msg == i18n.Model.NotFound {
			log.Logger.Infof("The Generator %d may have already been stopped", payload.Generator.ID)
			return nil
		}
		return fmt.Errorf("get generator failed: %s", ret.Msg)
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	ret = k8s.StopGenerator(ctx, generator)
	cancel()
	if !ret.OK {
		return fmt.Errorf("stop generator failed: %s", ret.Msg)
	}
	ret = db.InitGeneratorRepo(db.DB).Delete(generator.ID)
	if !ret.OK {
		return fmt.Errorf("delete generator failed: %s", ret.Msg)
	}
	return nil
}
