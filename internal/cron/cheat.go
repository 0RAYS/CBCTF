package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"time"
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
	for _, contest := range contests {
		if time.Now().Sub(contest.Start.Add(contest.Duration)) > job.Schedule*2 {
			continue
		}
		service.CheckWebReqIP(db.CronDB, contest)
		service.CheckVictimReqIP(db.CronDB, contest)
		service.CheckWrongFlag(db.CronDB, contest)
	}
	return model.SuccessRetVal()
}
