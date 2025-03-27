package service

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

func CalcSolversAndScore(tx *gorm.DB, flag model.Flag) (int64, float64, bool, string) {
	solvedCount, ok, msg := CountFlagSolved(tx, flag)
	if !ok {
		return 0, 0, false, msg
	}
	score := flag.CalcNewScore(solvedCount)
	if score != flag.CurrentScore {
		// 不考虑更新失败的情况, 不回滚
		repo := db.InitFlagRepo(tx)
		repo.Update(flag.ID, db.UpdateFlagOptions{CurrentScore: &score})
	}
	return solvedCount, score, true, "Success"
}
