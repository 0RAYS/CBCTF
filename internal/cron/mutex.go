package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"time"

	"github.com/robfig/cron/v3"
)

// clearSubmissionMutex 定时任务清理flag提交锁 service.SolvedMutex
func clearSubmissionMutex(c *cron.Cron) {
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(exec("ClearSubmissionMutex", func() error {
		contests := make(map[uint]model.Contest)
		contestRepo := db.InitContestRepo(db.DB)
		contestFlagRepo := db.InitContestFlagRepo(db.DB)
		service.SolvedMutex.Range(func(k, v any) bool {
			contestFlag, ret := contestFlagRepo.GetByID(k.(uint))
			if !ret.OK {
				service.SolvedMutex.Delete(k)
				return true
			}
			contest, ok := contests[contestFlag.ContestID]
			if !ok {
				contest, ret = contestRepo.GetByID(contestFlag.ContestID)
				if !ret.OK {
					service.SolvedMutex.Delete(k)
					return true
				}
				contests[contestFlag.ContestID] = contest
			}
			if !contest.IsRunning() {
				service.SolvedMutex.Delete(k)
			}
			return true
		})
		return nil
	})))
}

// clearCheatMutex 定时任务清理作弊检测锁 db.CheatMutex
func clearCheatMutex(c *cron.Cron) {
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(exec("ClearCheatMutex", func() error {
		contests := make(map[uint]model.Contest)
		contestRepo := db.InitContestRepo(db.DB)
		cheatRepo := db.InitCheatRepo(db.DB)
		db.CheatMutex.Range(func(k, v any) bool {
			hash := k.(string)
			cheat, ret := cheatRepo.Get(db.GetOptions{Conditions: map[string]any{"hash": hash}})
			if !ret.OK {
				db.CheatMutex.Delete(k)
				return true
			}
			contest, ok := contests[cheat.ContestID]
			if !ok {
				contest, ret = contestRepo.GetByID(cheat.ContestID)
				if !ret.OK {
					db.CheatMutex.Delete(k)
					return true
				}
				contests[cheat.ContestID] = contest
			}
			if !contest.IsRunning() {
				db.CheatMutex.Delete(k)
			}
			return true
		})
		return nil
	})))
}

// clearJoinTeamMutes 定时任务清理作弊检测锁 service.JoinTeamMutex
func clearJoinTeamMutes(c *cron.Cron) {
	c.Schedule(cron.Every(10*time.Minute), cron.FuncJob(exec("ClearJoinTeamMutes", func() error {
		contests := make(map[uint]model.Contest)
		contestRepo := db.InitContestRepo(db.DB)
		teamRepo := db.InitTeamRepo(db.DB)
		service.JoinTeamMutex.Range(func(k, v any) bool {
			contestFlag, ret := teamRepo.GetByID(k.(uint))
			if !ret.OK {
				service.JoinTeamMutex.Delete(k)
				return true
			}
			contest, ok := contests[contestFlag.ContestID]
			if !ok {
				contest, ret = contestRepo.GetByID(contestFlag.ContestID)
				if !ret.OK {
					service.JoinTeamMutex.Delete(k)
					return true
				}
				contests[contestFlag.ContestID] = contest
			}
			if contest.IsOver() {
				service.JoinTeamMutex.Delete(k)
			}
			return true
		})
		return nil
	})))
}
