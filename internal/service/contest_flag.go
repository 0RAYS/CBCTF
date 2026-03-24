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
		Preloads:   map[string]db.GetOptions{"TeamFlags": {Conditions: map[string]any{"team_id": team.ID}}},
	})
	if !ret.OK {
		return false, model.ContestFlag{}, model.TeamFlag{}, ret
	}
	if len(contestFlagL) == 0 {
		return false, model.ContestFlag{}, model.TeamFlag{}, model.RetVal{Msg: i18n.Model.ContestFlag.NotFound}
	}
	if contestChallenge.Type == model.QuestionChallengeType {
		contestFlag := contestFlagL[0]
		if len(contestFlag.TeamFlags) == 0 {
			return false, model.ContestFlag{}, model.TeamFlag{}, model.RetVal{Msg: i18n.Model.TeamFlag.NotFound}
		}
		teamFlag := contestFlag.TeamFlags[0]
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
