package service

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

func CalcContestFlagState(tx *gorm.DB, contestFlag model.ContestFlag) (int64, float64, bool, string) {
	solvers, ok, msg := db.InitSubmissionRepo(tx).CountWithConditions(db.GetOptions{
		{Key: "contest_flag_id", Value: contestFlag.ID, Op: "and"},
		{Key: "solved", Value: true, Op: "and"},
	})
	if !ok {
		return contestFlag.Solvers, contestFlag.CurrentScore, false, msg
	}
	return solvers, contestFlag.CalcScore(solvers - 1), true, i18n.Success
}
