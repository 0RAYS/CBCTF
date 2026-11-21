package service

import (
	"fmt"
	"sync"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
)

var SolvedMutex sync.Map

func Submit(
	tx *gorm.DB,
	user model.User,
	team model.Team,
	contest model.Contest,
	challenge model.Challenge,
	form dto.SubmitFlagForm,
	ip string,
) (model.Submission, model.RetVal) {
	submissionRepo := db.InitSubmissionRepo(tx)
	if challenge.Attempt != 0 {
		if ret := submissionRepo.LockAttemptScope(team.ID, challenge.ID); !ret.OK {
			return model.Submission{}, ret
		}
		attempts, ret := CountAttempts(tx, team, challenge)
		if !ret.OK {
			return model.Submission{}, ret
		}
		if challenge.Attempt <= attempts {
			return model.Submission{}, model.FailRetVal(i18n.Model.Submission.NotAllowed, "Submission", fmt.Sprintf("team_id=%d,challenge_id=%d,attempts=%d", team.ID, challenge.ID, attempts), nil)
		}
	}
	solved, challengeFlag, teamFlag, ret := VerifyFlag(tx, team, challenge, form.Flag)
	if !ret.OK {
		return model.Submission{}, ret
	}

	options := model.Submission{
		ContestID:       contest.ID,
		ChallengeID:     challenge.ID,
		TeamID:          team.ID,
		UserID:          user.ID,
		ChallengeFlagID: challengeFlag.ID,
		Value:           form.Flag,
		Score:           team.Score,
		Solved:          solved,
		IP:              ip,
	}

	if !solved {
		submission, createRet := submissionRepo.Create(options)
		if !createRet.OK {
			return model.Submission{}, createRet
		}
		prometheus.RecordFlagSubmission(contest.ID, string(challenge.Type), false)
		return submission, model.SuccessRetVal()
	}

	lockedFlag, ret := db.InitChallengeFlagRepo(tx).GetByIDForUpdate(challengeFlag.ID)
	if !ret.OK {
		return model.Submission{}, ret
	}

	lockedTeamFlag, ret := db.InitTeamFlagRepo(tx).GetByIDForUpdate(teamFlag.ID)
	if !ret.OK {
		return model.Submission{}, ret
	}
	if lockedTeamFlag.Solved {
		options.Solved = false
		submission, createRet := submissionRepo.Create(options)
		if !createRet.OK {
			return model.Submission{}, createRet
		}
		prometheus.RecordFlagSubmission(contest.ID, string(challenge.Type), false)
		return submission, model.SuccessRetVal()
	}

	existingSolvers, ret := submissionRepo.Count(db.CountOptions{
		Conditions: map[string]any{"challenge_flag_id": lockedFlag.ID, "solved": true},
	})
	if !ret.OK {
		return model.Submission{}, ret
	}

	submission, ret := submissionRepo.Create(options)
	if !ret.OK {
		return model.Submission{}, ret
	}

	if contest.Blood {
		switch existingSolvers {
		case 0:
			prometheus.RecordBlood(contest.ID, "first")
		case 1:
			prometheus.RecordBlood(contest.ID, "second")
		case 2:
			prometheus.RecordBlood(contest.ID, "third")
		}
	}

	teamFlagRepo := db.InitTeamFlagRepo(tx)
	if ret = teamFlagRepo.Update(lockedTeamFlag.ID, db.UpdateTeamFlagRepo{Solved: &solved}); !ret.OK {
		return model.Submission{}, ret
	}

	solvers, currentScore, ret := CalcChallengeFlagState(tx, lockedFlag)
	if !ret.OK {
		return model.Submission{}, ret
	}
	flagRepo := db.InitChallengeFlagRepo(tx)
	if ret = flagRepo.Update(lockedFlag.ID, db.UpdateChallengeFlagOptions{
		Solvers:      &solvers,
		CurrentScore: &currentScore,
		Last:         &submission.CreatedAt,
	}); !ret.OK {
		return model.Submission{}, ret
	}

	lockedTeam, ret := db.InitTeamRepo(tx).GetByIDForUpdate(team.ID)
	if !ret.OK {
		return model.Submission{}, ret
	}

	score, ret := CalcTeamScore(tx, contest, lockedTeam)
	if !ret.OK {
		return model.Submission{}, ret
	}
	teamRepo := db.InitTeamRepo(tx)
	if ret = teamRepo.Update(lockedTeam.ID, db.UpdateTeamOptions{
		Score: &score,
		Last:  &submission.CreatedAt,
	}); !ret.OK {
		return model.Submission{}, ret
	}
	if ret = submissionRepo.Update(submission.ID, db.UpdateSubmissionOptions{Score: &score}); !ret.OK {
		return model.Submission{}, ret
	}

	prometheus.RecordFlagSubmission(contest.ID, string(challenge.Type), true)
	return submission, model.SuccessRetVal()
}

func CountAttempts(tx *gorm.DB, team model.Team, challenge model.Challenge) (int64, model.RetVal) {
	return db.InitSubmissionRepo(tx).Count(db.CountOptions{
		Conditions: map[string]any{"team_id": team.ID, "challenge_id": challenge.ID, "solved": false},
	})
}

// CheckIfSolved 只有当前生效（active）的 flag 才参与"是否解出"的判定。
func CheckIfSolved(tx *gorm.DB, team model.Team, flags []model.ChallengeFlag) (bool, model.RetVal) {
	if len(flags) == 0 {
		return true, model.SuccessRetVal()
	}
	activeCount := 0
	for _, flag := range flags {
		if flag.Active {
			activeCount++
		}
	}
	solvedCount, ret := db.InitTeamFlagRepo(tx).CountSolvedForChallenge(team.ID, flags[0].ChallengeID)
	if !ret.OK {
		return false, ret
	}
	return solvedCount == int64(activeCount), model.SuccessRetVal()
}

func ListTeamSubmissions(tx *gorm.DB, team model.Team, form dto.ListModelsForm) ([]model.Submission, int64, model.RetVal) {
	return db.InitSubmissionRepo(tx).List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID},
		Sort:       []string{"id DESC"},
	})
}
