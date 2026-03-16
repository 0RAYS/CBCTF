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

const GenAttachmentTaskType = "tasks:attachment"

type GenAttachmentPayload struct {
	UserID    uint
	Generator model.Generator
	Challenge model.Challenge
	TeamID    uint
	Flags     []string
}

func EnqueueGenAttachmentTask(userID uint, generator model.Generator, challenge model.Challenge, team model.Team, teamFlags []model.TeamFlag) (*asynq.TaskInfo, error) {
	var flags []string
	for _, flag := range teamFlags {
		flags = append(flags, flag.Value)
	}
	payload, err := msgpack.Marshal(GenAttachmentPayload{userID, generator, challenge, team.ID, flags})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(GenAttachmentTaskType, payload)
	info, err := client.Enqueue(task, asynq.Queue(queueForTask(GenAttachmentTaskType)), asynq.MaxRetry(0), asynq.Timeout(2*time.Minute))
	if err == nil {
		prometheus.RecordTaskEnqueued(GenAttachmentTaskType)
	}
	return info, err
}

func HandleGenAttachmentTask(_ context.Context, t *asynq.Task) error {
	var payload GenAttachmentPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	ret := k8s.GenAttachment(ctx, payload.Challenge, payload.Generator, payload.TeamID, payload.Flags)
	cancel()
	db.InitGeneratorRepo(db.DB).UpdateStatus(payload.Generator.ID, ret.OK, time.Now())
	if !ret.OK {
		return fmt.Errorf("generate attachment failed: %s", ret.Msg)
	}
	return nil
}
