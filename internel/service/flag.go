package service

import (
	"CBCTF/internel/model"
	"gorm.io/gorm"
)

func CalcSolversAndScore(tx *gorm.DB, flag model.Flag) (int64, float64, bool, string) {
	solvedCount, ok, msg := CountFlagSolved(tx, flag)
	if !ok {
		return 0, 0, false, msg
	}
	score := flag.CalcNewScore(solvedCount)
	return solvedCount, score, true, "Success"
}
