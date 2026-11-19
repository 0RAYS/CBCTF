package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/task"
	"CBCTF/internal/utils"
)

// resolveGeneratorSource 按 rand_id 解析生成器来源。
// 比赛内使用比赛题目副本，题库题目只在测试模式下使用。
func resolveGeneratorSource(tx *gorm.DB, contestID uint, randID string) (model.GeneratorSourceSnapshot, model.RetVal) {
	if contestID > 0 {
		contestChallenge, ret := db.InitContestChallengeRepo(tx).GetByRandID(randID, db.GetOptions{
			Conditions: map[string]any{"contest_id": contestID},
		})
		if !ret.OK {
			return model.GeneratorSourceSnapshot{}, ret
		}
		if contestChallenge.Type != model.DynamicChallengeType {
			return model.GeneratorSourceSnapshot{}, model.RetVal{Msg: i18n.Model.Challenge.InvalidType}
		}
		return model.NewContestChallengeSource(contestChallenge), model.SuccessRetVal()
	}
	challenge, ret := db.InitChallengeRepo(tx).GetByRandID(randID)
	if !ret.OK {
		return model.GeneratorSourceSnapshot{}, ret
	}
	if challenge.Type != model.DynamicChallengeType {
		return model.GeneratorSourceSnapshot{}, model.RetVal{Msg: i18n.Model.Challenge.InvalidType}
	}
	return model.NewChallengeSource(challenge), model.SuccessRetVal()
}

func StartGenerators(tx *gorm.DB, contestID uint, form dto.StartGeneratorsForm) model.RetVal {
	if len(form.Challenges) == 0 {
		return model.SuccessRetVal()
	}
	generatorRepo := db.InitGeneratorRepo(tx)
	batch := model.NewBatch(len(form.Challenges))
	for index, randID := range form.Challenges {
		if tx.Statement.Context.Err() != nil {
			return batch.Result(tx.Statement.Context)
		}
		key := fmt.Sprintf("challenge:%s/instance:%d", randID, index+1)
		source, ret := resolveGeneratorSource(tx, contestID, randID)
		if !ret.OK {
			batch.Fail(key, "lookup", ret)
			if model.BatchDependencyFailed(ret) {
				return batch.Result(tx.Statement.Context)
			}
			continue
		}
		generator, ret := generatorRepo.Create(model.Generator{
			WorkerToken:        utils.RandHexStr(64),
			Image:              source.SourceImage(),
			ChallengeID:        source.SourceChallengeID(),
			ChallengeName:      source.SourceName(),
			ContestChallengeID: sql.Null[uint]{V: source.SourceContestChallengeID(), Valid: source.SourceContestChallengeID() > 0},
			ContestID:          sql.Null[uint]{V: source.SourceContestID(), Valid: source.SourceContestID() > 0},
			Name:               fmt.Sprintf("gen-%d-%d-%s", source.SourceContestID(), source.SourceID(), utils.RandHexStr(6)),
			Status:             model.WaitingGeneratorStatus,
		})
		if !ret.OK {
			batch.Fail(key, "create", ret)
			if model.BatchDependencyFailed(ret) {
				return batch.Result(tx.Statement.Context)
			}
			continue
		}
		if _, err := task.EnqueueStartGeneratorTask(source, generator); err != nil {
			batch.Fail(fmt.Sprintf("%s/generator:%d", key, generator.ID), "enqueue", model.RetVal{Msg: i18n.Task.EnqueueError})
			log.Logger.Warningf("Failed to enqueue start generator task: generator_id=%d name=%s source_id=%d error=%v", generator.ID, generator.Name, source.SourceID(), err)
			if ret := generatorRepo.Delete(generator.ID); !ret.OK {
				log.Logger.Warningf(
					"Failed to delete generator after enqueue failure: generator_id=%d name=%s source_id=%d reason=%s",
					generator.ID, generator.Name, source.SourceID(), ret.Msg,
				)
			}
			return batch.Result(tx.Statement.Context)
		}
		batch.Success(fmt.Sprintf("%s/generator:%d", key, generator.ID), "queued")
	}
	return batch.Result(tx.Statement.Context)
}

func warmGeneratorPool(source model.GeneratorSourceSnapshot) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := task.EnsureGeneratorPool(ctx, source); err != nil {
		log.Logger.Warningf("Failed to warm generator pool: source_id=%d error=%v", source.SourceID(), err)
	}
}

func StopGenerators(tx *gorm.DB, contestID uint, form dto.StopGeneratorsForm) model.RetVal {
	if len(form.Generators) == 0 {
		return model.SuccessRetVal()
	}
	options := db.GetOptions{Conditions: map[string]any{"id": form.Generators}, Deleted: true}
	if contestID > 0 {
		options.Conditions["contest_id"] = contestID
	}
	generators, _, ret := db.InitGeneratorRepo(tx).List(-1, -1, options)
	if !ret.OK {
		return ret
	}
	byID := make(map[uint]model.Generator, len(generators))
	for _, generator := range generators {
		byID[generator.ID] = generator
	}
	batch := model.NewBatch(len(form.Generators))
	for _, id := range form.Generators {
		if tx.Statement.Context.Err() != nil {
			return batch.Result(tx.Statement.Context)
		}
		key := fmt.Sprint(id)
		generator, exists := byID[id]
		if !exists {
			batch.Fail(key, "lookup", model.RetVal{Msg: i18n.Model.NotFound})
			continue
		}
		if generator.Status == model.TerminatingGeneratorStatus || generator.Status == model.StoppedGeneratorStatus {
			batch.Skip(key, "already_stopping")
			continue
		}
		ret = StopGenerator(tx, generator)
		if ret.OK {
			batch.Success(key, "queued")
		} else {
			batch.Fail(key, "stop", ret)
			if model.BatchDependencyFailed(ret) {
				return batch.Result(tx.Statement.Context)
			}
		}
	}
	return batch.Result(tx.Statement.Context)
}

func StopGenerator(tx *gorm.DB, generator model.Generator) model.RetVal {
	switch generator.Status {
	case model.WaitingGeneratorStatus, model.PendingGeneratorStatus:
		return model.RetVal{Msg: i18n.Model.Generator.NotStoppable}
	case model.TerminatingGeneratorStatus:
		return model.SuccessRetVal()
	}
	repo := db.InitGeneratorRepo(tx)
	if ret := repo.UpdateIfStatus(generator.ID, generator.Status, db.UpdateGeneratorOptions{Status: new(model.TerminatingGeneratorStatus)}); !ret.OK {
		return ret
	}
	generator.Status = model.TerminatingGeneratorStatus
	if err := unregisterGenerator(generator); err != nil {
		log.Logger.Warningf("Failed to unregister generator before stop: generator_id=%d name=%s error=%v", generator.ID, generator.Name, err)
	}
	err := task.EnqueueStopGeneratorTask(generator)
	if err != nil {
		log.Logger.Warningf("Failed to enqueue stop generator task: generator_id=%d name=%s challenge_id=%d error=%v", generator.ID, generator.Name, generator.ChallengeID, err)
		_ = repo.UpdateIfStatus(generator.ID, model.TerminatingGeneratorStatus, db.UpdateGeneratorOptions{Status: new(model.RunningGeneratorStatus)})
		generator.Status = model.RunningGeneratorStatus
		if registerErr := registerGenerator(generator); registerErr != nil {
			log.Logger.Warningf("Failed to re-register generator after stop enqueue failure: generator_id=%d name=%s error=%v", generator.ID, generator.Name, registerErr)
		}
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	log.Logger.Infof("Stop generator queued: generator_id=%d name=%s challenge_id=%d", generator.ID, generator.Name, generator.ChallengeID)
	return model.SuccessRetVal()
}

func stopGeneratorResources(tx *gorm.DB, options db.GetOptions) model.RetVal {
	generators, _, ret := db.InitGeneratorRepo(tx).List(-1, -1, options)
	if !ret.OK {
		return ret
	}
	batch := model.NewBatch(len(generators))
	for _, generator := range generators {
		if tx.Statement.Context.Err() != nil {
			return batch.Result(tx.Statement.Context)
		}
		key := fmt.Sprint(generator.ID)
		switch generator.Status {
		case model.WaitingGeneratorStatus, model.StoppedGeneratorStatus:
			batch.Skip(key, "not_provisioned")
			continue
		}
		steps := model.NewBatch(2)
		if err := unregisterGenerator(generator); err != nil {
			steps.Fail(key, "unregister", model.RetVal{Msg: i18n.Redis.DeleteError})
			log.Logger.Warningf("Failed to unregister generator before resource deletion: generator_id=%d name=%s error=%v", generator.ID, generator.Name, err)
		} else {
			steps.Success(key, "unregister")
		}
		ctx, cancel := context.WithTimeout(tx.Statement.Context, time.Minute)
		ret = k8s.StopGenerator(ctx, generator)
		cancel()
		steps.Record(key, "stop_resources", ret)
		batch.Record(key, "cleanup", steps.Result(tx.Statement.Context))
		if ret.OK {
			log.Logger.Infof("Stopped generator resources before model deletion: generator_id=%d name=%s challenge_id=%d", generator.ID, generator.Name, generator.ChallengeID)
		}
	}
	return batch.Result(tx.Statement.Context)
}

func GetGenerator(tx *gorm.DB, source model.GeneratorSource) (model.Generator, string, model.RetVal) {
	ctx, cancel := context.WithTimeout(tx.Statement.Context, 5*time.Second)
	defer cancel()
	target := source.SourceTarget()
	generator, lockToken, err := redis.LockAvailableGenerator(ctx, target)
	if err == nil {
		return generator, lockToken, model.SuccessRetVal()
	}
	if errors.Is(err, redis.ErrGeneratorPoolEmpty) {
		options := db.GetOptions{Conditions: map[string]any{
			"status": model.RunningGeneratorStatus,
		}}
		if target.ContestChallengeID > 0 {
			options.Conditions["contest_challenge_id"] = target.ContestChallengeID
			options.Conditions["contest_id"] = target.ContestID
		} else {
			options.Conditions["challenge_id"] = target.ChallengeID
			options.Conditions["contest_id"] = nil
		}
		generators, _, ret := db.InitGeneratorRepo(tx).List(-1, -1, options)
		if !ret.OK {
			return model.Generator{}, "", ret
		}
		if len(generators) == 0 {
			return model.Generator{}, "", model.RetVal{Msg: i18n.Model.Generator.NotAvailable}
		}
		for _, generator := range generators {
			if err := registerGenerator(generator); err != nil {
				return model.Generator{}, "", model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
			}
		}
		generator, lockToken, err = redis.LockAvailableGenerator(ctx, target)
		if err == nil {
			return generator, lockToken, model.SuccessRetVal()
		}
	}
	if errors.Is(err, redis.ErrGeneratorPoolEmpty) || errors.Is(err, redis.ErrNoAvailableGenerator) {
		return model.Generator{}, "", model.RetVal{Msg: i18n.Model.Generator.NotAvailable}
	}
	return model.Generator{}, "", model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
}

func registerGenerator(generator model.Generator) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return redis.RegisterGenerator(ctx, generator)
}

func unregisterGenerator(generator model.Generator) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return redis.UnregisterGenerator(ctx, generator)
}

func ListGenerators(tx *gorm.DB, contest model.Contest, form dto.ListGeneratorsForm) ([]model.Generator, int64, model.RetVal) {
	options := db.GetOptions{
		Deleted: form.Deleted,
		Sort:    []string{"id DESC"},
	}
	if contest.ID > 0 {
		options.Conditions = map[string]any{"contest_id": contest.ID}
	}
	return db.InitGeneratorRepo(tx).List(form.Limit, form.Offset, options)
}
