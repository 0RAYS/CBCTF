package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"time"
)

func checkCheatTask() model.RetVal {
	contests, _, ret := db.InitContestRepo(db.DB).List(-1, -1)
	if !ret.OK {
		return ret
	}
	for _, contest := range contests {
		if time.Now().Sub(contest.Start.Add(contest.Duration)) > 15*time.Minute {
			continue
		}
		service.CheckWebReqIP(db.DB, contest)
		service.CheckVictimReqIP(db.DB, contest)
		service.CheckWrongFlag(db.DB, contest)
		service.CheckSameDevice(db.DB, contest)
	}
	return model.SuccessRetVal()
}
