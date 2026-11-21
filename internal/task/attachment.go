package task

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
)

const genAttachmentTaskType = "tasks:attachment"

// GenAttachmentPayload 附件生成任务的载荷，题目以快照形式携带。
type GenAttachmentPayload struct {
	Flags  []string
	Source model.GeneratorSourceSnapshot
	UserID uint
	TeamID uint
}

func EnqueueGenAttachmentTask(userID uint, source model.GeneratorSourceSnapshot, teamID uint, teamFlags []model.TeamFlag) error {
	teamFlags = slices.Clone(teamFlags)
	slices.SortFunc(teamFlags, func(a, b model.TeamFlag) int {
		if a.ChallengeFlagID < b.ChallengeFlagID {
			return -1
		}
		if a.ChallengeFlagID > b.ChallengeFlagID {
			return 1
		}
		return 0
	})
	var flags []string
	for _, flag := range teamFlags {
		flags = append(flags, flag.Value)
	}
	payload, err := msgpack.Marshal(GenAttachmentPayload{
		UserID: userID,
		Source: source,
		TeamID: teamID,
		Flags:  flags,
	})
	if err != nil {
		return err
	}
	path := source.SourceAttachmentCachePath(teamID, flags)
	if info, err := os.Stat(path); err == nil && info.Size() > 0 {
		return nil
	}
	id := fmt.Sprintf("attachment-%x", sha256.Sum256([]byte(path)))
	task := asynq.NewTask(genAttachmentTaskType, payload)
	_, err = enqueueTask(
		genAttachmentTaskType, task,
		asynq.TaskID(id),
		asynq.MaxRetry(300),
		asynq.Deadline(time.Now().Add(30*time.Minute)),
		asynq.Timeout(90*time.Second),
	)
	if errors.Is(err, asynq.ErrTaskIDConflict) {
		info, inspectErr := inspector.GetTaskInfo(genAttachmentTaskType, id)
		if inspectErr != nil {
			return inspectErr
		}
		if info.State == asynq.TaskStateArchived {
			if err := inspector.DeleteTask(genAttachmentTaskType, id); err != nil {
				return err
			}
			_, err = enqueueTask(
				genAttachmentTaskType, task,
				asynq.TaskID(id),
				asynq.MaxRetry(300),
				asynq.Deadline(time.Now().Add(30*time.Minute)),
				asynq.Timeout(90*time.Second),
			)
			if errors.Is(err, asynq.ErrTaskIDConflict) {
				return nil
			}
			return err
		}
		return nil
	}
	return err
}

func HandleGenAttachmentTask(ctx context.Context, t *asynq.Task) error {
	payload, err := decodeMsgpackTask[GenAttachmentPayload](t)
	if err != nil {
		return err
	}
	path := payload.Source.SourceAttachmentCachePath(payload.TeamID, payload.Flags)
	if info, err := os.Stat(path); err == nil && info.Size() > 0 {
		return nil
	}
	generator, lockToken, err := redis.LockAvailableGenerator(ctx, payload.Source.SourceTarget())
	if err != nil {
		if errors.Is(err, redis.ErrGeneratorPoolEmpty) {
			if poolErr := EnsureGeneratorPool(ctx, payload.Source); poolErr != nil {
				return poolErr
			}
		}
		return err
	}
	defer unlockGeneratorAttachment(generator.ID, lockToken)
	current, ret := db.InitGeneratorRepo(db.TaskDB.WithContext(ctx)).GetByID(generator.ID)
	if !ret.OK {
		return taskResourceError("get attachment generator", ret)
	}
	if current.Status != model.RunningGeneratorStatus {
		if err := redis.UnregisterGenerator(ctx, current); err != nil {
			return err
		}
		return redis.ErrNoAvailableGenerator
	}
	if current.Image != payload.Source.SourceImage() {
		if err := redis.UnregisterGenerator(ctx, current); err != nil {
			return err
		}
		if err := EnqueueStopGeneratorTask(current); err != nil {
			return err
		}
		if err := EnsureGeneratorPool(ctx, payload.Source); err != nil {
			return err
		}
		return redis.ErrNoAvailableGenerator
	}
	generator = current

	generateCtx, cancel := context.WithTimeout(ctx, 80*time.Second)
	ret = k8s.GenAttachment(generateCtx, payload.Source, generator, payload.TeamID, payload.Flags)
	cancel()
	// Persist the outcome and release broken pool members even after cancellation.
	finishCtx, finishCancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer finishCancel()
	generatorRepo := db.InitGeneratorRepo(db.TaskDB.WithContext(finishCtx))
	if statusRet := generatorRepo.UpdateStatus(generator.ID, ret.OK, time.Now()); !statusRet.OK {
		log.Logger.Warningf("Failed to record attachment outcome: generator_id=%d reason=%v", generator.ID, statusRet)
	}
	if !ret.OK {
		if ret.Msg == i18n.Model.NotFound || ret.Msg == i18n.K8S.NotFound {
			unregisterCtx, unregisterCancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := redis.UnregisterGenerator(unregisterCtx, generator); err != nil {
				log.Logger.Warningf("Failed to unregister generator: generator_id=%d error=%v", generator.ID, err)
			}
			unregisterCancel()
			if deleteRet := generatorRepo.Delete(generator.ID); !deleteRet.OK {
				return fmt.Errorf("generate attachment failed: %w; delete unavailable generator failed: %w", ret.AsError(), deleteRet.AsError())
			}
		}
		return taskResourceError("generate attachment failed", ret)
	}
	log.Logger.Infof(
		"Attachment generated: user_id=%d team_id=%d source_id=%d generator_id=%d",
		payload.UserID, payload.TeamID, payload.Source.SourceID(), generator.ID,
	)
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
