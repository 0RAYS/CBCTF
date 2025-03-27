package service

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

// IsSolved model.Usage 需要预加载
func IsSolved(tx *gorm.DB, team model.Team, usage model.Usage) bool {
	var (
		count                 int
		submissionRepo        = db.InitSubmissionRepo(tx)
		submissions, _, ok, _ = submissionRepo.GetAllByKeyID("team_id", team.ID, -1, -1, false, 0, true)
	)
	if !ok {
		return false
	}
	for _, submission := range submissions {
		if submission.UsageID == usage.ID {
			count++
		}
	}
	switch usage.Challenge.Type {
	case model.StaticChallenge, model.DynamicChallenge, model.DockerChallenge:
		if count < 1 {
			return false
		}
	case model.DockersChallenge:
		if count != len(usage.Flags) {
			return false
		}
	default:
		return false
	}
	return true
}

// CountAttempts 统计题目的尝试次数
func CountAttempts(tx *gorm.DB, team model.Team, usage model.Usage) int64 {
	var count int64
	submissionRepo := db.InitSubmissionRepo(tx)
	submissions, _, ok, _ := submissionRepo.GetAllByKeyID("team_id", team.ID, -1, -1, false, 0, false)
	if !ok {
		return count
	}
	for _, submission := range submissions {
		if submission.UsageID == usage.ID {
			count++
		}
	}
	return count
}

// CountFlagSolved 统计指定 model.Flag 的解题次数
func CountFlagSolved(tx *gorm.DB, flag model.Flag) (int64, bool, string) {
	var (
		count                   int64
		submissionRepo          = db.InitSubmissionRepo(tx)
		submissions, _, ok, msg = submissionRepo.GetAllByKeyID("contest_id", flag.ContestID, -1, -1, true, 0, true)
	)
	if !ok {
		return count, false, msg
	}
	for _, submission := range submissions {
		if submission.FlagID == flag.ID {
			count++
		}
	}
	if count < flag.Solvers {
		// 不考虑更新失败的情况, 不回滚
		flagRepo := db.InitFlagRepo(tx)
		flagRepo.Update(flag.ID, db.UpdateFlagOptions{Solvers: &count})
	}
	return count, true, "Success"
}
