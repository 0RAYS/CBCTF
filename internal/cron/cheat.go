package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"time"
)

func checkCheatTask() model.RetVal {
	job, ret := db.InitCronJobRepo(db.DB).GetByUniqueField("name", model.ClearEmptyTeamCronJob)
	if !ret.OK {
		return ret
	}
	contests, _, ret := db.InitContestRepo(db.DB).List(-1, -1)
	if !ret.OK {
		return ret
	}
	for _, contest := range contests {
		if time.Now().Sub(contest.Start.Add(contest.Duration)) > job.Schedule*2 {
			continue
		}
		service.CheckWebReqIP(db.DB, contest)
		service.CheckVictimReqIP(db.DB, contest)
		service.CheckWrongFlag(db.DB, contest)
		service.CheckSameDevice(db.DB, contest)
	}
	return model.SuccessRetVal()
}
