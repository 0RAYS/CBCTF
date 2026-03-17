package task

import (
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

func HandleStartGeneratorTask(ctx context.Context, t *asynq.Task) error {
	var payload StartGeneratorPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	challenge := payload.Challenge
	generator := payload.Generator
	generatorRepo := db.InitGeneratorRepo(db.DB)
	ret := generatorRepo.Update(generator.ID, db.UpdateGeneratorOptions{Status: new(model.PendingGeneratorStatus)})
	if !ret.OK {
		return fmt.Errorf("start generator fail, update generator fail: %s", ret.Msg)
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	_, ret = k8s.StartGenerator(ctx, challenge, generator)
	cancel()
	if !ret.OK {
		_, err := EnqueueStopGeneratorTask(generator)
		return err
	}
	ret = generatorRepo.Update(generator.ID, db.UpdateGeneratorOptions{Status: new(model.RunningGeneratorStatus)})
	if !ret.OK {
		return fmt.Errorf("start generator fail, update generator fail: %s", ret.Msg)
	}
	return nil
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
	generator := payload.Generator
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	ret := k8s.StopGenerator(ctx, generator)
	cancel()
	if !ret.OK {
		return fmt.Errorf("stop generator fail: %s", ret.Msg)
	}
	ret = db.InitGeneratorRepo(db.DB).Delete(generator.ID)
	if !ret.OK {
		return fmt.Errorf("stop generator fail, delete generator fail: %s", ret.Msg)
	}
	return nil
}
