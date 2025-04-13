package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

func VerifyFlag(tx *gorm.DB, team model.Team, usage model.Usage, value string) (bool, model.Flag, model.Answer, bool, string) {
	flagRepo := db.InitFlagRepo(tx)
	flags, _, ok, msg := flagRepo.GetByKeyID("usage_id", usage.ID, -1, -1, "Answers")
	if !ok {
		return false, model.Flag{}, model.Answer{}, false, msg
	}
	for _, flag := range flags {
		for _, answer := range flag.Answers {
			if answer.TeamID == team.ID && answer.Value == value {
				if answer.Solved {
					return false, flag, answer, true, "AlreadySolved"
				}
				return true, flag, answer, true, "Success"
			}
		}
	}
	// 没有找到答案, 则默认为第一个flag
	return false, flags[0], model.Answer{}, true, "FlagNotMatch"
}

func UpdateFlag(tx *gorm.DB, flag model.Flag, form f.UpdateFlagForm) (bool, string) {
	flagRepo := db.InitFlagRepo(tx)
	challengeRepo := db.InitChallengeRepo(tx)
	challenge, ok, msg := challengeRepo.GetByID(flag.Usage.ChallengeID)
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
	case model.DynamicChallenge, model.PodChallenge, model.PodsChallenge:
		options.Value = form.Value
	default:
		return false, "InvalidChallengeType"
	}
	return flagRepo.Update(flag.ID, options)
}

func CalcSolversAndScore(tx *gorm.DB, flag model.Flag) (int64, float64, bool, string) {
	count, ok, msg := db.InitSubmissionRepo(tx).CountByKeyID("flag_id", flag.ID, true)
	if !ok {
		return 0, 0, false, msg
	}
	score := flag.CalcCurrentScore(count - 1)
	return count, score, true, "Success"
}
