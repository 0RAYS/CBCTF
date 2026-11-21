package task

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
)

// EnsureGeneratorPool Pool admission is serialized across platform replicas; capacity is reserved
// by waiting/pending records before a start task is enqueued.
// 生成器池按题目 id 划分，不同题目的池互不干扰。
func EnsureGeneratorPool(ctx context.Context, source model.GeneratorSourceSnapshot) error {
	if source.Type != model.DynamicChallengeType || config.Env.K8S.GeneratorPoolSize <= 0 {
		return nil
	}
	// 锁 id 必须唯一标识一个生成器池，题目 id 全局唯一。
	lockID := uint64(source.ChallengeID)
	return db.WithWorkloadLock(ctx, db.WorkloadLockDB, "generator-pool", uint(lockID), func() error {
		repo := db.InitGeneratorRepo(db.TaskDB.WithContext(ctx))
		conditions := map[string]any{
			"image": source.Image,
			"status": []string{
				model.WaitingGeneratorStatus,
				model.PendingGeneratorStatus,
				model.RunningGeneratorStatus,
			},
		}
		conditions["challenge_id"] = source.ChallengeID
		generators, _, ret := repo.List(-1, -1, db.GetOptions{Conditions: conditions})
		if !ret.OK {
			return taskResourceError("list generator pool", ret)
		}
		capacity := 0
		for _, generator := range generators {
			if generator.Status == model.PendingGeneratorStatus && time.Since(generator.UpdatedAt) > 3*time.Minute {
				if err := EnqueueStopGeneratorTask(generator); err != nil {
					return err
				}
				if ret := repo.UpdateIfStatus(generator.ID, model.PendingGeneratorStatus, db.UpdateGeneratorOptions{
					Status: new(model.TerminatingGeneratorStatus),
				}); !ret.OK {
					return taskResourceError("retire stalled generator", ret)
				}
				continue
			}
			capacity++
			if generator.Status == model.RunningGeneratorStatus {
				if err := redis.RegisterGenerator(ctx, generator); err != nil {
					return err
				}
			}
		}
		for i := capacity; i < config.Env.K8S.GeneratorPoolSize; i++ {
			generator, ret := repo.Create(model.Generator{
				ChallengeID:   source.ChallengeID,
				ChallengeName: source.Name,
				ContestID:     sql.Null[uint]{V: source.ContestID, Valid: source.ContestID > 0},
				Name:          fmt.Sprintf("gen-%d-%d-%s", source.ContestID, source.SourceID(), utils.RandHexStr(12)),
				WorkerToken:   utils.RandHexStr(64),
				Image:         source.Image,
				Status:        model.WaitingGeneratorStatus,
			})
			if !ret.OK {
				return taskResourceError("create generator pool member", ret)
			}
			if _, err := EnqueueStartGeneratorTask(source, generator); err != nil {
				_ = repo.Delete(generator.ID)
				return err
			}
		}
		return nil
	})
}
