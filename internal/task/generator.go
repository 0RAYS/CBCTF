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
	return enqueueTask(startGeneratorTaskType, task, asynq.MaxRetry(0), asynq.Timeout(3*time.Minute))
}

func HandleStartGeneratorTask(ctx context.Context, t *asynq.Task) error {
	var payload StartGeneratorPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return db.WithWorkloadLock(ctx, db.WorkloadLockDB, "generator", payload.Generator.ID, func() error {
		cleanupQueued := false
		claimed := false
		err := func() error {
			challenge := payload.Challenge
			generator := payload.Generator
			generatorRepo := db.InitGeneratorRepo(db.TaskDB)
			currentGenerator, ret := generatorRepo.GetByID(generator.ID)
			if !ret.OK {
				if ret.Msg == i18n.Model.NotFound {
					log.Logger.Debugf("Start generator skipped: generator_id=%d no longer exists", generator.ID)
					return nil
				}
				return fmt.Errorf("get generator failed: %s", ret.Msg)
			}
			if currentGenerator.Status != model.WaitingGeneratorStatus {
				log.Logger.Infof("Start generator skipped: generator_id=%d status=%s", generator.ID, currentGenerator.Status)
				return nil
			}
			if ret = generatorRepo.UpdateIfStatus(generator.ID, model.WaitingGeneratorStatus, db.UpdateGeneratorOptions{Status: new(model.PendingGeneratorStatus)}); !ret.OK {
				return fmt.Errorf("update generator failed: %s", ret.Msg)
			}
			claimed = true
			generator = currentGenerator
			startCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
			log.Logger.Infof("Starting generator provisioning: generator_id=%d name=%s challenge_id=%d", generator.ID, generator.Name, challenge.ID)
			_, ret = k8s.StartGenerator(startCtx, challenge, generator)
			cancel()
			if !ret.OK {
				_ = generatorRepo.UpdateIfStatus(generator.ID, model.PendingGeneratorStatus, db.UpdateGeneratorOptions{Status: new(model.TerminatingGeneratorStatus)})
				if err := EnqueueStopGeneratorTask(generator); err != nil {
					return fmt.Errorf("%w; enqueue cleanup failed: %v", taskResourceError("start generator failed", ret), err)
				}
				cleanupQueued = true
				return taskResourceError("start generator failed", ret)
			}
			ret = generatorRepo.UpdateIfStatus(generator.ID, model.PendingGeneratorStatus, db.UpdateGeneratorOptions{Status: new(model.RunningGeneratorStatus)})
			if !ret.OK {
				return fmt.Errorf("update generator failed: %s", ret.Msg)
			}
			if err := ctx.Err(); err != nil {
				return err
			}
			generator.Status = model.RunningGeneratorStatus
			registerCtx, registerCancel := context.WithTimeout(ctx, 5*time.Second)
			if err := redis.RegisterGenerator(registerCtx, generator); err != nil {
				registerCancel()
				return fmt.Errorf("register generator failed: %w", err)
			}
			registerCancel()
			log.Logger.Infof("Generator is running: generator_id=%d name=%s challenge_id=%d", generator.ID, generator.Name, challenge.ID)
			return nil
		}()
		if err != nil && claimed && !cleanupQueued {
			repo := db.InitGeneratorRepo(db.TaskDB)
			for _, status := range []string{model.PendingGeneratorStatus, model.RunningGeneratorStatus} {
				_ = repo.UpdateIfStatus(payload.Generator.ID, status, db.UpdateGeneratorOptions{Status: new(model.TerminatingGeneratorStatus)})
			}
			if enqueueErr := EnqueueStopGeneratorTask(payload.Generator); enqueueErr != nil {
				cleanupCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
				_ = redis.UnregisterGenerator(cleanupCtx, payload.Generator)
				ret := k8s.StopGenerator(cleanupCtx, payload.Generator)
				cancel()
				if ret.OK {
					_ = repo.Delete(payload.Generator.ID)
				}
				log.Logger.Warningf("Failed to enqueue generator cleanup after start failure: generator_id=%d error=%v", payload.Generator.ID, enqueueErr)
			}
		}
		return err
	})
}

type StopGeneratorPayload struct {
	Generator model.Generator
}

func EnqueueStopGeneratorTask(generator model.Generator) error {
	payload, err := msgpack.Marshal(StopGeneratorPayload{generator})
	if err != nil {
		return err
	}
	task := asynq.NewTask(stopGeneratorTaskType, payload)
	_, err = enqueueTask(stopGeneratorTaskType, task, asynq.MaxRetry(3), asynq.Timeout(2*time.Minute))
	return err
}

func HandleStopGeneratorTask(ctx context.Context, t *asynq.Task) error {
	var payload StopGeneratorPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return db.WithWorkloadLock(ctx, db.WorkloadLockDB, "generator", payload.Generator.ID, func() error {
		generatorRepo := db.InitGeneratorRepo(db.TaskDB)
		generator, ret := generatorRepo.GetByID(payload.Generator.ID)
		if !ret.OK {
			if ret.Msg == i18n.Model.NotFound {
				log.Logger.Debugf("Stop generator skipped: generator_id=%d no longer exists", payload.Generator.ID)
				return nil
			}
			return fmt.Errorf("get generator failed: %s", ret.Msg)
		}
		log.Logger.Infof("Stopping generator: generator_id=%d name=%s challenge_id=%d", generator.ID, generator.Name, generator.ChallengeID)
		unregisterCtx, unregisterCancel := context.WithTimeout(ctx, 5*time.Second)
		if err := redis.UnregisterGenerator(unregisterCtx, generator); err != nil {
			log.Logger.Warningf("Failed to unregister generator before stop: generator_id=%d error=%v", generator.ID, err)
		}
		unregisterCancel()
		ctx, cancel := context.WithTimeout(ctx, time.Minute)
		ret = k8s.StopGenerator(ctx, generator)
		cancel()
		if !ret.OK {
			return fmt.Errorf("stop generator failed: %s", ret.Msg)
		}
		ret = db.InitGeneratorRepo(db.TaskDB).Delete(generator.ID)
		if !ret.OK {
			return fmt.Errorf("delete generator failed: %s", ret.Msg)
		}
		log.Logger.Infof("Generator stopped: generator_id=%d name=%s challenge_id=%d", generator.ID, generator.Name, generator.ChallengeID)
		return nil
	})
}
