package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

func UpdateFlag(tx *gorm.DB, flag model.Flag, form f.UpdateFlagForm) (bool, string) {
	flagRepo := db.InitFlagRepo(tx)
	challengeRepo := db.InitChallengeRepo(tx)
	challenge, ok, msg := challengeRepo.GetByID(flag.Usage.ChallengeID, true, 0)
	if !ok {
		return ok, msg
	}
	options := db.UpdateFlagOptions{
		Score:     form.Score,
		Decay:     form.Decay,
		MinScore:  form.MinScore,
		ScoreType: form.ScoreType,
		Attempt:   form.Attempt,
	}
	switch challenge.Type {
	case model.StaticChallenge:
		options.Value = &flag.Value
	case model.DynamicChallenge, model.DockerChallenge, model.DockersChallenge:
		options.Value = form.Value
	default:
		return false, "InvalidChallengeType"
	}
	return flagRepo.Update(flag.ID, options)
}

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
