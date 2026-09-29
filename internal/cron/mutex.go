package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
)

// clearSubmissionMutexTask 定时任务清理flag提交锁 service.SolvedMutex
func clearSubmissionMutexTask() model.RetVal {
	result := model.SuccessRetVal()
	contests := make(map[uint]model.Contest)
	contestRepo := db.InitContestRepo(db.CronDB)
	contestFlagRepo := db.InitContestFlagRepo(db.CronDB)
	service.SolvedMutex.Range(func(k, _ any) bool {
		contestFlag, ret := contestFlagRepo.GetByID(k.(uint))
		if !ret.OK {
			if ret.Msg != i18n.Model.NotFound {
				result = ret
				return false
			}
			service.SolvedMutex.Delete(k)
			return true
		}
		contest, ok := contests[contestFlag.ContestID]
		if !ok {
			contest, ret = contestRepo.GetByID(contestFlag.ContestID)
			if !ret.OK {
				if ret.Msg != i18n.Model.NotFound {
					result = ret
					return false
				}
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
	return result
}

// clearCheatMutexTask 定时任务清理作弊检测锁 db.CheatMutex
func clearCheatMutexTask() model.RetVal {
	result := model.SuccessRetVal()
	contests := make(map[uint]model.Contest)
	contestRepo := db.InitContestRepo(db.CronDB)
	cheatRepo := db.InitCheatRepo(db.CronDB)
	db.CheatMutex.Range(func(k, _ any) bool {
		hash := k.(string)
		cheat, ret := cheatRepo.Get(db.GetOptions{Conditions: map[string]any{"hash": hash}})
		if !ret.OK {
			if ret.Msg != i18n.Model.NotFound {
				result = ret
				return false
			}
			db.CheatMutex.Delete(k)
			return true
		}
		contest, ok := contests[cheat.ContestID]
		if !ok {
			contest, ret = contestRepo.GetByID(cheat.ContestID)
			if !ret.OK {
				if ret.Msg != i18n.Model.NotFound {
					result = ret
					return false
				}
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
	return result
}

// clearJoinTeamMutexTask 定时任务清理加入队伍锁 service.JoinTeamMutex
func clearJoinTeamMutexTask() model.RetVal {
	result := model.SuccessRetVal()
	contests := make(map[uint]model.Contest)
	contestRepo := db.InitContestRepo(db.CronDB)
	teamRepo := db.InitTeamRepo(db.CronDB)
	service.JoinTeamMutex.Range(func(k, _ any) bool {
		contestFlag, ret := teamRepo.GetByID(k.(uint))
		if !ret.OK {
			if ret.Msg != i18n.Model.NotFound {
				result = ret
				return false
			}
			service.JoinTeamMutex.Delete(k)
			return true
		}
		contest, ok := contests[contestFlag.ContestID]
		if !ok {
			contest, ret = contestRepo.GetByID(contestFlag.ContestID)
			if !ret.OK {
				if ret.Msg != i18n.Model.NotFound {
					result = ret
					return false
				}
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
	return result
}
