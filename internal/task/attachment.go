package task

import (
	"CBCTF/internal/db"
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

const genAttachmentTaskType = "tasks:attachment"

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
	task := asynq.NewTask(genAttachmentTaskType, payload)
	info, err := client.Enqueue(task, asynq.Queue(genAttachmentTaskType), asynq.MaxRetry(0), asynq.Timeout(2*time.Minute))
	if err == nil {
		prometheus.RecordTaskEnqueued(genAttachmentTaskType)
	}
	return info, err
}

func HandleGenAttachmentTask(ctx context.Context, t *asynq.Task) error {
	var payload GenAttachmentPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	log.Logger.Infof("Generating attachment: user_id=%d team_id=%d challenge_id=%d generator_id=%d", payload.UserID, payload.TeamID, payload.Challenge.ID, payload.Generator.ID)
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	ret := k8s.GenAttachment(ctx, payload.Challenge, payload.Generator, payload.TeamID, payload.Flags)
	cancel()
	db.InitGeneratorRepo(db.DB).UpdateStatus(payload.Generator.ID, ret.OK, time.Now())
	if !ret.OK {
		return fmt.Errorf("generate attachment failed: %s", ret.Msg)
	}
	log.Logger.Infof("Attachment generated: user_id=%d team_id=%d challenge_id=%d generator_id=%d", payload.UserID, payload.TeamID, payload.Challenge.ID, payload.Generator.ID)
	return nil
}
