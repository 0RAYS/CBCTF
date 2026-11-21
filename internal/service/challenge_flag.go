package service

import (
	"fmt"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
)

// VerifyFlag 校验队伍提交的 flag 是否命中该题目的生效 flag。
func VerifyFlag(tx *gorm.DB, team model.Team, challenge model.Challenge, value string) (bool, model.ChallengeFlag, model.TeamFlag, model.RetVal) {
	flagL, _, ret := db.InitChallengeFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"challenge_id": challenge.ID},
	})
	if !ret.OK {
		return false, model.ChallengeFlag{}, model.TeamFlag{}, ret
	}
	activeFlagL := make([]model.ChallengeFlag, 0, len(flagL))
	for _, flag := range flagL {
		if flag.Active {
			activeFlagL = append(activeFlagL, flag)
		}
	}
	teamFlag, ret := db.InitTeamFlagRepo(tx).GetByChallengeAndValue(team.ID, challenge.ID, value)
	if !ret.OK && ret.Msg != i18n.Model.TeamFlag.NotFound {
		return false, model.ChallengeFlag{}, model.TeamFlag{}, ret
	}
	if ret.OK {
		for _, flag := range activeFlagL {
			if flag.ID != teamFlag.ChallengeFlagID {
				continue
			}
			if teamFlag.Solved {
				return false, flag, teamFlag, model.RetVal{OK: true, Msg: i18n.Model.TeamFlag.AlreadySolved}
			}
			return true, flag, teamFlag, model.SuccessRetVal()
		}
	}
	if len(activeFlagL) == 0 {
		return false, model.ChallengeFlag{}, model.TeamFlag{}, model.FailRetVal(i18n.Model.ChallengeFlag.NotFound, "ChallengeFlag", fmt.Sprintf("challenge_id=%d", challenge.ID), nil)
	}
	// 没有找到答案, 则默认为第一个flag
	return false, activeFlagL[0], model.TeamFlag{}, model.RetVal{OK: true, Msg: i18n.Model.TeamFlag.NotMatch}
}

func CalcChallengeFlagState(tx *gorm.DB, flag model.ChallengeFlag) (int64, float64, model.RetVal) {
	solvers, ret := db.InitSubmissionRepo(tx).Count(db.CountOptions{
		Conditions: map[string]any{"challenge_flag_id": flag.ID, "solved": true},
	})
	if !ret.OK {
		return 0, 0, ret
	}
	return solvers, flag.CalcScore(solvers - 1), model.SuccessRetVal()
}

func SubmitChallengeFlag(
	tx *gorm.DB,
	user model.User,
	team model.Team,
	contest model.Contest,
	challenge model.Challenge,
	form dto.SubmitFlagForm,
	ip string,
) model.RetVal {
	var solved bool
	ret := db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		_, submitRet := Submit(tx2, user, team, contest, challenge, form, ip)
		if !submitRet.OK {
			return submitRet
		}

		flags, _, listRet := db.InitChallengeFlagRepo(tx2).List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"challenge_id": challenge.ID},
		})
		if !listRet.OK {
			return listRet
		}
		if challenge.Type == model.PodsChallengeType {
			var checkRet model.RetVal
			solved, checkRet = CheckIfSolved(tx2, team, flags)
			if !checkRet.OK {
				return checkRet
			}
		}
		return model.SuccessRetVal()
	})
	if !ret.OK {
		return ret
	}
	if solved {
		queueTeamVictimStop(tx, team.ID, challenge.ID)
	}
	return model.SuccessRetVal()
}

func ListChallengeFlags(tx *gorm.DB, challenge model.Challenge) ([]model.ChallengeFlag, model.RetVal) {
	flags, _, ret := db.InitChallengeFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"challenge_id": challenge.ID},
	})
	return flags, ret
}
