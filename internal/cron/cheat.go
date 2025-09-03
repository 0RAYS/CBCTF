package cron

import (
	"CBCTF/internal/cron/cheat"
	"CBCTF/internal/db"
	"time"

	"github.com/robfig/cron/v3"
)

func checkCheat(c *cron.Cron) {
	function := exec("CheckCheat", func() {
		contests, _, ok, _ := db.InitContestRepo(db.DB).List(-1, -1, db.GetOptions{
			Selects: []string{"id", "start", "duration"},
		})
		if !ok {
			return
		}
		for _, contest := range contests {
			if time.Now().Sub(contest.Start.Add(contest.Duration)) > 15*time.Minute {
				continue
			}
			cheat.CheckVictimReqIP(contest)
			cheat.CheckWrongFlag(contest)
			cheat.CheckSameDevice(contest)
		}
	})
	function()
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(function))
}
