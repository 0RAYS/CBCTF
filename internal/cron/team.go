package cron

import (
	"fmt"
	"time"

	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

func clearEmptyTeamTask() model.RetVal {
	job, ret := db.InitCronJobRepo(db.CronDB).GetByUniqueField("name", model.ClearEmptyTeamCronJob)
	if !ret.OK {
		return ret
	}
	contests, _, ret := db.InitContestRepo(db.CronDB).List(-1, -1)
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
	repo := db.InitTeamRepo(db.CronDB)
	teams, _, ret := repo.List(-1, -1, db.GetOptions{Conditions: map[string]any{"contest_id": contestIDL}})
	if !ret.OK {
		return ret
	}
	teamIDL := make([]uint, 0, len(teams))
	for _, team := range teams {
		teamIDL = append(teamIDL, team.ID)
	}
	userCountMap, ret := repo.CountUsersMap(teamIDL...)
	if !ret.OK {
		return ret
	}
	batch := model.NewBatch(len(teams))
	for _, team := range teams {
		if db.CronDB.Statement.Context.Err() != nil {
			return batch.Result(db.CronDB.Statement.Context)
		}
		key := fmt.Sprint(team.ID)
		if userCountMap[team.ID] != 0 {
			batch.Skip(key, "not_empty")
			continue
		}
		ret = db.WithTransactionDB(db.CronDB, func(tx *db.Tx) model.RetVal { return db.InitTeamRepo(tx).Delete(team.ID) })
		if ret.OK {
			batch.Success(key, "deleted")
			log.Logger.Infof("Delete empty team: %d", team.ID)
		} else {
			batch.Fail(key, "delete", ret)
			if model.BatchDependencyFailed(ret) {
				return batch.Result(db.CronDB.Statement.Context)
			}
		}
	}
	return batch.Result(db.CronDB.Statement.Context)
}
