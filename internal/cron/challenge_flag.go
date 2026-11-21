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
	contestIDs := make([]uint, 0, len(contests))
	for _, contest := range contests {
		if contest.EndedWithin(time.Now(), job.Schedule*2) {
			contestIDs = append(contestIDs, contest.ID)
		}
	}
	challengeIDL := make([]uint, 0)
	if len(contestIDs) > 0 {
		challenges, _, ret := db.InitChallengeRepo(root).List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"contest_id": contestIDs},
		})
		if !ret.OK {
			return ret
		}
		for _, challenge := range challenges {
			challengeIDL = append(challengeIDL, challenge.ID)
		}
	}
	flags, _, ret := db.InitChallengeFlagRepo(root).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"challenge_id": challengeIDL},
	})
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
			repo := db.InitChallengeFlagRepo(tx)
			current, ret := repo.GetByIDForUpdate(flag.ID)
			if !ret.OK {
				return ret
			}
			solvers, score, ret := service.CalcChallengeFlagState(tx, current)
			if !ret.OK {
				return ret
			}
			if solvers == current.Solvers && score == current.CurrentScore {
				return model.SuccessRetVal()
			}
			return repo.Update(current.ID, db.UpdateChallengeFlagOptions{CurrentScore: &score, Solvers: &solvers})
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
