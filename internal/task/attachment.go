package task

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
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
	LockToken string
}

func EnqueueGenAttachmentTask(userID uint, generator model.Generator, lockToken string, challenge model.Challenge, team model.Team, teamFlags []model.TeamFlag) (*asynq.TaskInfo, error) {
	var flags []string
	for _, flag := range teamFlags {
		flags = append(flags, flag.Value)
	}
	payload, err := msgpack.Marshal(GenAttachmentPayload{
		UserID:    userID,
		Generator: generator,
		Challenge: challenge,
		TeamID:    team.ID,
		Flags:     flags,
		LockToken: lockToken,
	})
	if err != nil {
		unlockGeneratorAttachment(generator.ID, lockToken)
		return nil, err
	}
	task := asynq.NewTask(genAttachmentTaskType, payload)
	info, err := enqueueTask(genAttachmentTaskType, task, asynq.MaxRetry(0), asynq.Timeout(5*time.Minute))
	if err != nil {
		unlockGeneratorAttachment(generator.ID, lockToken)
	}
	return info, err
}

func HandleGenAttachmentTask(ctx context.Context, t *asynq.Task) error {
	var payload GenAttachmentPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	log.Logger.Infof("Generating attachment: user_id=%d team_id=%d challenge_id=%d generator_id=%d", payload.UserID, payload.TeamID, payload.Challenge.ID, payload.Generator.ID)
	lockToken := payload.LockToken
	if lockToken != "" {
		valid, err := redis.RefreshGeneratorAttachmentLock(ctx, payload.Generator.ID, lockToken)
		if err != nil {
			db.InitGeneratorRepo(db.TaskDB).UpdateStatus(payload.Generator.ID, false, time.Now())
			return err
		}
		if !valid {
			lockToken = ""
		}
	}
	if lockToken == "" {
		var err error
		lockToken, err = redis.LockGeneratorAttachment(ctx, payload.Generator.ID)
		if err != nil {
			db.InitGeneratorRepo(db.TaskDB).UpdateStatus(payload.Generator.ID, false, time.Now())
			return err
		}
	}
	defer unlockGeneratorAttachment(payload.Generator.ID, lockToken)

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	ret := k8s.GenAttachment(ctx, payload.Challenge, payload.Generator, payload.TeamID, payload.Flags)
	cancel()
	generatorRepo := db.InitGeneratorRepo(db.TaskDB)
	generatorRepo.UpdateStatus(payload.Generator.ID, ret.OK, time.Now())
	if !ret.OK {
		if ret.Msg == i18n.Model.NotFound || ret.Msg == i18n.K8S.NotFound {
			unregisterCtx, unregisterCancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := redis.UnregisterGenerator(unregisterCtx, payload.Generator); err != nil {
				log.Logger.Warningf("Failed to unregister generator: generator_id=%d error=%v", payload.Generator.ID, err)
			}
			unregisterCancel()
			if deleteRet := generatorRepo.Delete(payload.Generator.ID); !deleteRet.OK {
				return fmt.Errorf("generate attachment failed: %s; delete unavailable generator failed: %s", ret.Msg, deleteRet.Msg)
			}
		}
		return fmt.Errorf("generate attachment failed: %s", ret.Msg)
	}
	log.Logger.Infof("Attachment generated: user_id=%d team_id=%d challenge_id=%d generator_id=%d", payload.UserID, payload.TeamID, payload.Challenge.ID, payload.Generator.ID)
	return nil
}

func unlockGeneratorAttachment(generatorID uint, lockToken string) {
	if lockToken == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redis.UnlockGeneratorAttachment(ctx, generatorID, lockToken); err != nil {
		log.Logger.Warningf("Failed to unlock generator attachment: generator_id=%d error=%v", generatorID, err)
	}
}
