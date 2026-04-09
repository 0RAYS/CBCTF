package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"slices"
	"strings"

	"gorm.io/gorm"
)

func VerifyFlag(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge, value string) (bool, model.ContestFlag, model.TeamFlag, model.RetVal) {
	contestFlagRepo := db.InitContestFlagRepo(tx)
	contestFlagL, _, ret := contestFlagRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		return false, model.ContestFlag{}, model.TeamFlag{}, ret
	}
	if len(contestFlagL) == 0 {
		return false, model.ContestFlag{}, model.TeamFlag{}, model.RetVal{Msg: i18n.Model.ContestFlag.NotFound}
	}
	if contestChallenge.Type == model.QuestionChallengeType {
		contestFlag := contestFlagL[0]
		teamFlag, ret := db.InitTeamFlagRepo(tx).Get(db.GetOptions{
			Conditions: map[string]any{"team_id": team.ID, "contest_flag_id": contestFlag.ID},
		})
		if !ret.OK {
			return false, model.ContestFlag{}, model.TeamFlag{}, model.RetVal{Msg: i18n.Model.TeamFlag.NotFound}
		}
		optionsIDL := strings.Split(contestFlag.Value, ",")
		answerIDL := strings.Split(value, ",")
		if len(optionsIDL) != len(answerIDL) {
			return false, contestFlag, model.TeamFlag{}, model.RetVal{OK: true, Msg: i18n.Model.TeamFlag.NotMatch}
		}
		for _, answerID := range answerIDL {
			if !slices.Contains(optionsIDL, answerID) {
				return false, contestFlag, model.TeamFlag{}, model.RetVal{OK: true, Msg: i18n.Model.TeamFlag.NotMatch}
			}
		}
		if teamFlag.Solved {
			return false, contestFlag, teamFlag, model.RetVal{OK: true, Msg: i18n.Model.TeamFlag.AlreadySolved}
		}
		return true, contestFlag, teamFlag, model.SuccessRetVal()
	}

	teamFlag, ret := db.InitTeamFlagRepo(tx).GetByContestChallengeAndValue(team.ID, contestChallenge.ID, value)
	if ret.OK {
		for _, contestFlag := range contestFlagL {
			if contestFlag.ID != teamFlag.ContestFlagID {
				continue
			}
			if teamFlag.Solved {
				return false, contestFlag, teamFlag, model.RetVal{OK: true, Msg: i18n.Model.TeamFlag.AlreadySolved}
			}
			return true, contestFlag, teamFlag, model.SuccessRetVal()
		}
	}
	// 没有找到答案, 则默认为第一个flag
	return false, contestFlagL[0], model.TeamFlag{}, model.RetVal{OK: true, Msg: i18n.Model.TeamFlag.NotMatch}
}

func CalcContestFlagState(tx *gorm.DB, contestFlag model.ContestFlag) (int64, float64, model.RetVal) {
	solvers, ret := db.InitSubmissionRepo(tx).Count(db.CountOptions{
		Conditions: map[string]any{"contest_flag_id": contestFlag.ID, "solved": true},
	})
	if !ret.OK {
		return 0, 0, ret
	}
	return solvers, contestFlag.CalcScore(solvers - 1), model.SuccessRetVal()
}

func SubmitContestFlag(tx *gorm.DB, user model.User, team model.Team, contest model.Contest, challenge model.Challenge, contestChallenge model.ContestChallenge, form dto.SubmitFlagForm, ip string) model.RetVal {
	var solved bool
	ret := db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		_, submitRet := Submit(tx2, user, team, contest, contestChallenge, form, ip)
		if !submitRet.OK {
			return submitRet
		}

		contestFlags, _, listRet := db.InitContestFlagRepo(tx2).List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
		})
		if !listRet.OK {
			return listRet
		}
		solved = contestChallenge.Type == model.PodsChallengeType && CheckIfSolved(tx2, team, contestFlags)
		return model.SuccessRetVal()
	})
	if !ret.OK {
		return ret
	}
	if solved {
		go func() {
			victim, victimRet := db.InitVictimRepo(tx).HasAliveVictim(team.ID, challenge.ID)
			if !victimRet.OK {
				return
			}
			_ = ForceStopVictim(tx, victim)
		}()
	}
	return model.SuccessRetVal()
}

func ListContestFlags(tx *gorm.DB, contestChallenge model.ContestChallenge) ([]model.ContestFlag, model.RetVal) {
	flags, _, ret := db.InitContestFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	return flags, ret
}

func UpdateContestFlag(tx *gorm.DB, contestChallenge model.ContestChallenge, contestFlag model.ContestFlag, form dto.UpdateContestFlagForm) model.RetVal {
	if contestChallenge.Type == model.QuestionChallengeType && form.Value != nil {
		form.Value = &contestFlag.Value
	}
	currentScore := contestFlag.CurrentScore
	if form.Score != nil && *form.Score < currentScore {
		currentScore = *form.Score
	}
	return db.InitContestFlagRepo(tx).Update(contestFlag.ID, db.UpdateContestFlagOptions{
		Value:        form.Value,
		Score:        form.Score,
		CurrentScore: &currentScore,
		Decay:        form.Decay,
		MinScore:     form.MinScore,
		ScoreType:    form.ScoreType,
	})
}
