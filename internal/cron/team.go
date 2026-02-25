package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"errors"
	"time"

	"github.com/robfig/cron/v3"
)

func ClearEmptyTeam(c *cron.Cron) {
	function := exec("ClearEmptyTeam", func() error {
		contests, _, ret := db.InitContestRepo(db.DB).List(-1, -1)
		if !ret.OK {
			return errors.New(ret.Msg)
		}
		contestIDL := make([]uint, 0)
		for _, contest := range contests {
			if time.Now().Sub(contest.Start.Add(contest.Duration)) > 10*time.Minute {
				continue
			}
			contestIDL = append(contestIDL, contest.ID)
		}
		repo := db.InitTeamRepo(db.DB)
		teams, _, ret := repo.List(-1, -1, db.GetOptions{Conditions: map[string]any{"contest_id": contestIDL}})
		if !ret.OK {
			return errors.New(ret.Msg)
		}
		for _, team := range teams {
			if repo.CountAssociation(team, "Users") == 0 {
				if ret = repo.Delete(team.ID); ret.OK {
					log.Logger.Infof("Delete empty team: %d", team.ID)
				}
			}
		}
		return nil
	})
	function()
	c.Schedule(cron.Every(5*time.Minute), cron.FuncJob(function))
}
