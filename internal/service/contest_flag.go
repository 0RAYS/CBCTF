package service

import (
	"CBCTF/internal/db"
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
		Preloads:   map[string]db.GetOptions{"TeamFlags": {}},
	})
	if !ret.OK {
		return false, model.ContestFlag{}, model.TeamFlag{}, ret
	}
	if contestChallenge.Type == model.QuestionChallengeType {
		if len(contestFlagL) == 0 {
			return false, model.ContestFlag{}, model.TeamFlag{}, model.RetVal{Msg: i18n.Model.NotFound, Attr: map[string]any{"Model": model.ContestFlag{}.GetModelName()}}
		}
		optionsIDL := strings.Split(contestFlagL[0].Value, ",")
		answerIDL := strings.Split(value, ",")
		if len(optionsIDL) != len(answerIDL) {
			return false, contestFlagL[0], model.TeamFlag{}, model.RetVal{OK: true, Msg: i18n.Model.TeamFlag.NotMatch}
		}
		for _, answerID := range answerIDL {
			if !slices.Contains(optionsIDL, answerID) {
				return false, contestFlagL[0], model.TeamFlag{}, model.RetVal{OK: true, Msg: i18n.Model.TeamFlag.NotMatch}
			}
		}
		for _, teamFlag := range contestFlagL[0].TeamFlags {
			if teamFlag.TeamID == team.ID {
				if teamFlag.Solved {
					return false, contestFlagL[0], teamFlag, model.RetVal{OK: true, Msg: i18n.Model.TeamFlag.AlreadySolved}
				}
				return true, contestFlagL[0], teamFlag, model.SuccessRetVal()
			}
		}
		return false, contestFlagL[0], model.TeamFlag{}, model.RetVal{OK: true, Msg: i18n.Model.NotFound, Attr: map[string]any{"Model": model.TeamFlag{}.GetModelName()}}
	}
	for _, contestFlag := range contestFlagL {
		for _, teamFlag := range contestFlag.TeamFlags {
			if teamFlag.TeamID == team.ID && teamFlag.Value == value {
				if teamFlag.Solved {
					return false, contestFlag, teamFlag, model.RetVal{OK: true, Msg: i18n.Model.TeamFlag.AlreadySolved}
				}
				return true, contestFlag, teamFlag, model.SuccessRetVal()
			}
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
