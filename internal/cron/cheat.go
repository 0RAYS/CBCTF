package cron

import (
	"fmt"
	"time"

	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
)

func checkCheatTask() model.RetVal {
	job, ret := db.InitCronJobRepo(db.CronDB).GetByUniqueField("name", model.CheckCheatCronJob)
	if !ret.OK {
		return ret
	}
	contests, _, ret := db.InitContestRepo(db.CronDB).List(-1, -1)
	if !ret.OK {
		return ret
	}
	batch := model.NewBatch(len(contests))
	for _, contest := range contests {
		if db.CronDB.Statement.Context.Err() != nil {
			return batch.Result(db.CronDB.Statement.Context)
		}
		key := fmt.Sprint(contest.ID)
		if time.Now().Sub(contest.Start.Add(contest.Duration)) > job.Schedule*2 {
			batch.Skip(key, "outside_scan_window")
			continue
		}
		ret := service.RunCheatChecks(db.CronDB, contest)
		batch.Record(key, "scan", ret)
		if !ret.OK && model.BatchDependencyFailed(ret) {
			return batch.Result(db.CronDB.Statement.Context)
		}
	}
	return batch.Result(db.CronDB.Statement.Context)
}
