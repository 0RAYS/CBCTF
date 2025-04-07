package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

func VerifyFlag(tx *gorm.DB, team model.Team, usage model.Usage, value string) (bool, model.Flag, bool, string) {
	flagRepo := db.InitFlagRepo(tx)
	flags, _, ok, msg := flagRepo.GetByKeyID("usage_id", usage.ID, -1, -1, true, 0)
	if !ok {
		return false, model.Flag{}, false, msg
	}
	for _, flag := range flags {
		for _, answer := range flag.Answers {
			if answer.TeamID == team.ID && answer.Value == value {
				if answer.Solved {
					return true, flag, false, "AlreadySolved"
				}
				return true, flag, true, "Success"
			}
		}
	}
	// 没有找到答案, 则默认为第一个flag
	return false, flags[0], false, "Success"
}

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
