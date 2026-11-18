package service

import (
	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/view"
)

// ListChallengeUsages 返回引用了指定题库题目的比赛题目实例。
//
// 比赛题目在加入时已深拷贝为独立副本，题库的改动不会影响已加入的比赛，
// 因此这里只是提示影响范围，不阻断题库题目的修改与删除。
func ListChallengeUsages(tx *gorm.DB, challenge model.Challenge) ([]view.ChallengeUsageView, model.RetVal) {
	usages, ret := db.InitChallengeRepo(tx).ListContestChallengeUsages(challenge.ID)
	if !ret.OK {
		return nil, ret
	}
	views := make([]view.ChallengeUsageView, 0, len(usages))
	for _, usage := range usages {
		usageView := view.ChallengeUsageView{
			ContestID:        usage.ContestID,
			ContestName:      usage.Contest.Name,
			ContestChallenge: usage,
		}
		if count, countRet := db.InitTeamFlagRepo(tx).CountSolvedTeamsForContestChallenge(usage.ID); countRet.OK {
			usageView.Solvers = count
		}
		views = append(views, usageView)
	}
	return views, model.SuccessRetVal()
}

// CountChallengeUsages 返回引用了指定题库题目的比赛题目实例数量。
func CountChallengeUsages(tx *gorm.DB, challenge model.Challenge) (int64, model.RetVal) {
	usages, ret := db.InitChallengeRepo(tx).ListContestChallengeUsages(challenge.ID)
	if !ret.OK {
		return 0, ret
	}
	return int64(len(usages)), model.SuccessRetVal()
}
