package cron

import (
	"fmt"
	"time"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"CBCTF/internal/task"
)

func warmChallengeImagesTask() model.RetVal {
	challenges, ret := db.InitChallengeRepo(db.CronDB).ListForUnfinishedContests(time.Now())
	if !ret.OK {
		return ret
	}
	batch := model.NewBatch(len(challenges))
	for _, challenge := range challenges {
		if db.CronDB.Statement.Context.Err() != nil {
			return batch.Result(db.CronDB.Statement.Context)
		}
		if err := task.EnqueuePrepullTask(service.ChallengeImages(challenge)); err != nil {
			log.Logger.Warningf("Scheduled warmup enqueue failed: challenge_id=%d error=%v", challenge.ID, err)
			batch.Fail(fmt.Sprint(challenge.ID), "enqueue", model.RetVal{Msg: i18n.Task.EnqueueError})
			return batch.Result(db.CronDB.Statement.Context)
		}
		batch.Success(fmt.Sprint(challenge.ID), "queued")
	}
	return batch.Result(db.CronDB.Statement.Context)
}
