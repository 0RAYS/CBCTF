package service

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

func VerifyFlag(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge, value string) (bool, model.ContestFlag, model.TeamFlag, bool, string) {
	contestFlagRepo := db.InitContestFlagRepo(tx)
	contestFlagL, _, ok, msg := contestFlagRepo.ListWithConditions(-1, -1, db.GetOptions{
		{Key: "contest_challenge_id", Value: contestChallenge.ID, Op: "and"},
	}, false, "TeamFlags")
	if !ok {
		return false, model.ContestFlag{}, model.TeamFlag{}, false, msg
	}
	for _, contestFlag := range contestFlagL {
		for _, teamFlag := range contestFlag.TeamFlags {
			if teamFlag.TeamID == team.ID && teamFlag.Value == value {
				if teamFlag.Solved {
					return false, contestFlag, teamFlag, true, i18n.AlreadySolved
				}
				return true, contestFlag, teamFlag, true, i18n.Success
			}
		}
	}
	// 没有找到答案, 则默认为第一个flag
	return false, contestFlagL[0], model.TeamFlag{}, true, i18n.FlagNotMatch
}

func CalcContestFlagState(tx *gorm.DB, contestFlag model.ContestFlag) (int64, float64, bool, string) {
	solvers, ok, msg := db.InitSubmissionRepo(tx).CountWithConditions(db.GetOptions{
		{Key: "contest_flag_id", Value: contestFlag.ID, Op: "and"},
		{Key: "solved", Value: true, Op: "and"},
	}, false)
	if !ok {
		return contestFlag.Solvers, contestFlag.CurrentScore, false, msg
	}
	return solvers, contestFlag.CalcScore(solvers - 1), true, i18n.Success
}
