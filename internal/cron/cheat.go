package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/service"
	"errors"
	"time"

	"github.com/robfig/cron/v3"
)

func checkCheat(c *cron.Cron) {
	function := exec("CheckCheat", func() error {
		contests, _, ret := db.InitContestRepo(db.DB).List(-1, -1)
		if !ret.OK {
			return errors.New(ret.Msg)
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
		return nil
	})
	function()
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(function))
}
