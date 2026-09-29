package cron

import (
	"context"
	"fmt"
	"time"

	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
)

func updateFlagScoreTask() model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	root := db.CronDB.WithContext(ctx)
	job, ret := db.InitCronJobRepo(root).GetByUniqueField("name", model.UpdateFlagScoreCronJob)
	if !ret.OK {
		return ret
	}
	contests, ret := db.InitContestRepo(root).FindAll()
	if !ret.OK {
		return ret
	}
	ids := make([]uint, 0, len(contests))
	for _, contest := range contests {
		if time.Since(contest.Start.Add(contest.Duration)) <= job.Schedule*2 {
			ids = append(ids, contest.ID)
		}
	}
	flags, ret := db.InitContestFlagRepo(root).FindAll(db.GetOptions{Conditions: map[string]any{"contest_id": ids}})
	if !ret.OK {
		return ret
	}
	batch := model.NewBatch(len(flags))
	for _, flag := range flags {
		if ctx.Err() != nil {
			return batch.Result(ctx)
		}
		ret = db.WithTransactionDB(root, func(tx *db.Tx) model.RetVal {
			// Share the database row lock with flag submission across replicas.
			repo := db.InitContestFlagRepo(tx)
			current, ret := repo.GetByIDForUpdate(flag.ID)
			if !ret.OK {
				return ret
			}
			solvers, score, ret := service.CalcContestFlagState(tx, current)
			if !ret.OK {
				return ret
			}
			if solvers == current.Solvers && score == current.CurrentScore {
				return model.SuccessRetVal()
			}
			return repo.Update(current.ID, db.UpdateContestFlagOptions{CurrentScore: &score, Solvers: &solvers})
		})
		key := fmt.Sprint(flag.ID)
		if ret.OK {
			batch.Success(key, "updated")
		} else {
			batch.Fail(key, "update_score", ret)
			if model.BatchDependencyFailed(ret) {
				return batch.Result(ctx)
			}
		}
	}
	return batch.Result(ctx)
}
