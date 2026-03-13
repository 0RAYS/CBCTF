package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
)

// clearSubmissionMutexTask 定时任务清理flag提交锁 service.SolvedMutex
func clearSubmissionMutexTask() model.RetVal {
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
	return model.SuccessRetVal()
}

// clearCheatMutexTask 定时任务清理作弊检测锁 db.CheatMutex
func clearCheatMutexTask() model.RetVal {
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
	return model.SuccessRetVal()
}

// clearJoinTeamMutexTask 定时任务清理加入队伍锁 service.JoinTeamMutex
func clearJoinTeamMutexTask() model.RetVal {
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
	return model.SuccessRetVal()
}
