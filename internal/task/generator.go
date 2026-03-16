package task

import (
	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/utils"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	StartGeneratorTaskType = "tasks:generator:start"
	StopGeneratorTaskType  = "tasks:generator:stop"
)

type StartGeneratorPayload struct {
	ContestID uint
	Challenge model.Challenge
}

func EnqueueStartGeneratorTask(contestID uint, challenge model.Challenge) (*asynq.TaskInfo, error) {
	payload, err := msgpack.Marshal(StartGeneratorPayload{contestID, challenge})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(StartGeneratorTaskType, payload)
	info, err := client.Enqueue(task, asynq.MaxRetry(0), asynq.Timeout(2*time.Minute))
	if err != nil {
		prometheus.RecordTaskEnqueued(StartGeneratorTaskType)
	}
	return info, err
}

func HandleStartGeneratorTask(_ context.Context, t *asynq.Task) error {
	var payload StartGeneratorPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	contestID := payload.ContestID
	challenge := payload.Challenge
	generatorRepo := db.InitGeneratorRepo(db.DB)
	generator, ret := generatorRepo.Create(db.CreateGeneratorOptions{
		ChallengeID: challenge.ID,
		ContestID:   sql.Null[uint]{V: contestID, Valid: contestID > 0},
		Name:        fmt.Sprintf("gen-%d-%d-%s", contestID, challenge.ID, utils.RandStr(6)),
	})
	if !ret.OK {
		return fmt.Errorf("start generator fail, create generator fail: %s", ret.Msg)
	}
	ret = generatorRepo.Update(generator.ID, db.UpdateGeneratorOptions{Status: new(model.PendingGeneratorStatus)})
	if !ret.OK {
		return fmt.Errorf("start generator fail, update generator fail: %s", ret.Msg)
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
	task := asynq.NewTask(StopGeneratorTaskType, payload)
	info, err := client.Enqueue(task, asynq.MaxRetry(0), asynq.Timeout(2*time.Minute))
	if err != nil {
		prometheus.RecordTaskEnqueued(StopGeneratorTaskType)
	}
	return info, err
}

func HandleStopGeneratorTask(_ context.Context, t *asynq.Task) error {
	var payload StopGeneratorPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	generator := payload.Generator
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
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
