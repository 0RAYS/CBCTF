package task

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/redis"
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	genAttachmentTaskType = "tasks:attachment"
)

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
	info, err := client.Enqueue(task, asynq.Queue(genAttachmentTaskType), asynq.MaxRetry(0), asynq.Timeout(5*time.Minute))
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
	lockToken, err := redis.LockGeneratorAttachment(ctx, payload.Generator.ID)
	if err != nil {
		db.InitGeneratorRepo(db.DB).UpdateStatus(payload.Generator.ID, false, time.Now())
		return err
	}
	defer func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = redis.UnlockGeneratorAttachment(unlockCtx, payload.Generator.ID, lockToken); err != nil {
			log.Logger.Warningf("Failed to unlock generator attachment: generator_id=%d error=%v", payload.Generator.ID, err)
		}
	}()
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	ret := k8s.GenAttachment(ctx, payload.Challenge, payload.Generator, payload.TeamID, payload.Flags)
	cancel()
	generatorRepo := db.InitGeneratorRepo(db.DB)
	_ = generatorRepo.UpdateStatus(payload.Generator.ID, ret.OK, time.Now())
	if !ret.OK {
		if ret.Msg == i18n.Model.NotFound || ret.Msg == i18n.K8S.NotFound {
			if deleteRet := generatorRepo.Delete(payload.Generator.ID); !deleteRet.OK {
				return fmt.Errorf("generate attachment failed: %s; delete unavailable generator failed: %s", ret.Msg, deleteRet.Msg)
			}
		}
		return fmt.Errorf("generate attachment failed: %s", ret.Msg)
	}
	log.Logger.Infof("Attachment generated: user_id=%d team_id=%d challenge_id=%d generator_id=%d", payload.UserID, payload.TeamID, payload.Challenge.ID, payload.Generator.ID)
	return nil
}
