package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"time"
)

func clearEmptyTeamTask() model.RetVal {
	job, ret := db.InitCronJobRepo(db.DB).GetByUniqueField("name", model.ClearEmptyTeamCronJob)
	if !ret.OK {
		return ret
	}
	contests, _, ret := db.InitContestRepo(db.DB).List(-1, -1)
	if !ret.OK {
		return ret
	}
	contestIDL := make([]uint, 0)
	for _, contest := range contests {
		if time.Now().Sub(contest.Start.Add(contest.Duration)) > job.Schedule*2 {
			continue
		}
		contestIDL = append(contestIDL, contest.ID)
	}
	repo := db.InitTeamRepo(db.DB)
	teams, _, ret := repo.List(-1, -1, db.GetOptions{Conditions: map[string]any{"contest_id": contestIDL}})
	if !ret.OK {
		return ret
	}
	for _, team := range teams {
		if repo.CountAssociation(team, "Users") == 0 {
			if ret = repo.Delete(team.ID); ret.OK {
				log.Logger.Infof("Delete empty team: %d", team.ID)
			}
		}
	}
	return model.SuccessRetVal()
}
