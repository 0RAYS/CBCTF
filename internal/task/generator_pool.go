package task

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"context"
	"database/sql"
	"fmt"
)

// Pool admission is serialized across platform replicas; capacity is reserved
// by waiting/pending records before a start task is enqueued.
func EnsureGeneratorPool(ctx context.Context, contestID uint, challenge model.Challenge) error {
	if challenge.Type != model.DynamicChallengeType || config.Env.K8S.GeneratorPoolSize <= 0 {
		return nil
	}
	return db.WithWorkloadLock(ctx, db.WorkloadLockDB, "generator-pool", challenge.ID, func() error {
		repo := db.InitGeneratorRepo(db.TaskDB.WithContext(ctx))
		conditions := map[string]any{"challenge_id": challenge.ID, "image": challenge.GeneratorImage, "status": []string{model.WaitingGeneratorStatus, model.PendingGeneratorStatus, model.RunningGeneratorStatus}, "contest_id": nil}
		if contestID > 0 {
			conditions["contest_id"] = contestID
		}
		generators, _, ret := repo.List(-1, -1, db.GetOptions{Conditions: conditions})
		if !ret.OK {
			return taskResourceError("list generator pool", ret)
		}
		for _, generator := range generators {
			if generator.Status == model.RunningGeneratorStatus {
				if err := redis.RegisterGenerator(ctx, generator); err != nil {
					return err
				}
			}
		}
		for i := len(generators); i < config.Env.K8S.GeneratorPoolSize; i++ {
			generator, ret := repo.Create(model.Generator{ChallengeID: challenge.ID, ChallengeName: challenge.Name, ContestID: sql.Null[uint]{V: contestID, Valid: contestID > 0}, Name: fmt.Sprintf("gen-%d-%d-%s", contestID, challenge.ID, utils.RandHexStr(12)), WorkerToken: utils.RandHexStr(64), Image: challenge.GeneratorImage, Status: model.WaitingGeneratorStatus})
			if !ret.OK {
				return taskResourceError("create generator pool member", ret)
			}
			if _, err := EnqueueStartGeneratorTask(challenge, generator); err != nil {
				_ = repo.Delete(generator.ID)
				return err
			}
		}
		return nil
	})
}
