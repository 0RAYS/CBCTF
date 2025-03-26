package service

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

func GetContestFlag(tx *gorm.DB, contestID uint) ([]model.Flag, bool, string) {
	var (
		repo              = db.InitFlagRepo(tx)
		flags, _, ok, msg = repo.GetByKeyID("contest_id", contestID, -1, -1, true, 3)
	)
	return flags, ok, msg
}

func CalcSolversAndScore(tx *gorm.DB, flag model.Flag) (int64, float64, bool, string) {
	solvedCount, ok, msg := CountFlagSolved(tx, flag)
	if !ok {
		return 0, 0, false, msg
	}
	score := flag.CalcNewScore(solvedCount)
	return solvedCount, score, true, "Success"
}
