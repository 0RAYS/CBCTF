package cron

import (
	"sync"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
)

// clearContestMutexes drops the lock entries whose record or contest is gone,
// and the ones whose contest no longer accepts that kind of lock.
func clearContestMutexes(
	mutex *sync.Map,
	contestIDOf func(key any) (uint, model.RetVal),
	expired func(model.Contest) bool,
) model.RetVal {
	result := model.SuccessRetVal()
	contests := make(map[uint]model.Contest)
	contestRepo := db.InitContestRepo(db.CronDB)
	mutex.Range(func(k, _ any) bool {
		contestID, ret := contestIDOf(k)
		if !ret.OK {
			if ret.Msg != i18n.Model.NotFound {
				result = ret
				return false
			}
			mutex.Delete(k)
			return true
		}
		contest, ok := contests[contestID]
		if !ok {
			contest, ret = contestRepo.GetByID(contestID)
			if !ret.OK {
				if ret.Msg != i18n.Model.NotFound {
					result = ret
					return false
				}
				mutex.Delete(k)
				return true
			}
			contests[contestID] = contest
		}
		if expired(contest) {
			mutex.Delete(k)
		}
		return true
	})
	return result
}

// clearSubmissionMutexTask 定时任务清理flag提交锁 service.SolvedMutex
func clearSubmissionMutexTask() model.RetVal {
	flagRepo := db.InitChallengeFlagRepo(db.CronDB)
	challengeRepo := db.InitChallengeRepo(db.CronDB)
	return clearContestMutexes(&service.SolvedMutex, func(key any) (uint, model.RetVal) {
		flag, ret := flagRepo.GetByID(key.(uint))
		if !ret.OK {
			return 0, ret
		}
		challenge, ret := challengeRepo.GetByID(flag.ChallengeID)
		if !ret.OK {
			return 0, ret
		}
		return challenge.ContestID.V, model.SuccessRetVal()
	}, func(contest model.Contest) bool {
		return !contest.IsRunning()
	})
}

// clearCheatMutexTask 定时任务清理作弊检测锁 db.CheatMutex
func clearCheatMutexTask() model.RetVal {
	cheatRepo := db.InitCheatRepo(db.CronDB)
	return clearContestMutexes(&db.CheatMutex, func(key any) (uint, model.RetVal) {
		cheat, ret := cheatRepo.Get(db.GetOptions{Conditions: map[string]any{"hash": key.(string)}})
		return cheat.ContestID, ret
	}, func(contest model.Contest) bool {
		return !contest.IsRunning()
	})
}

// clearJoinTeamMutexTask 定时任务清理加入队伍锁 service.JoinTeamMutex
func clearJoinTeamMutexTask() model.RetVal {
	teamRepo := db.InitTeamRepo(db.CronDB)
	return clearContestMutexes(&service.JoinTeamMutex, func(key any) (uint, model.RetVal) {
		team, ret := teamRepo.GetByID(key.(uint))
		return team.ContestID, ret
	}, func(contest model.Contest) bool {
		return contest.IsOver()
	})
}
