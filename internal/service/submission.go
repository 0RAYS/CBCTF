package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"sync"

	"gorm.io/gorm"
)

var SolvedMutex sync.Map

func Submit(tx *gorm.DB, user model.User, team model.Team, contest model.Contest, contestChallenge model.ContestChallenge, form dto.SubmitFlagForm, ip string) (model.Submission, model.RetVal) {
	if contestChallenge.Attempt != 0 && contestChallenge.Attempt <= CountAttempts(tx, team, contestChallenge) {
		return model.Submission{}, model.RetVal{Msg: i18n.Model.Submission.NotAllowed}
	}
	submissionRepo := db.InitSubmissionRepo(tx)
	solved, contestFlag, teamFlag, ret := VerifyFlag(tx, team, contestChallenge, form.Flag)
	if !ret.OK {
		return model.Submission{}, ret
	}

	options := db.CreateSubmissionOptions{
		ContestChallengeID: contestChallenge.ID,
		ContestID:          contest.ID,
		ChallengeID:        contestChallenge.ChallengeID,
		TeamID:             team.ID,
		UserID:             user.ID,
		ContestFlagID:      contestFlag.ID,
		Value:              form.Flag,
		Score:              team.Score,
		Solved:             solved,
		IP:                 ip,
	}

	if !solved {
		submission, ret := submissionRepo.Create(options)
		if !ret.OK {
			return model.Submission{}, ret
		}
		prometheus.RecordFlagSubmission(contest.ID, string(contestChallenge.Type), false)
		return submission, model.SuccessRetVal()
	}

	lockedContestFlag, ret := db.InitContestFlagRepo(tx).GetByIDForUpdate(contestFlag.ID)
	if !ret.OK {
		return model.Submission{}, ret
	}

	lockedTeamFlag, ret := db.InitTeamFlagRepo(tx).GetByIDForUpdate(teamFlag.ID)
	if !ret.OK {
		return model.Submission{}, ret
	}
	if lockedTeamFlag.Solved {
		options.Solved = false
		submission, ret := submissionRepo.Create(options)
		if !ret.OK {
			return model.Submission{}, ret
		}
		prometheus.RecordFlagSubmission(contest.ID, string(contestChallenge.Type), false)
		return submission, model.SuccessRetVal()
	}

	existingSolvers, ret := submissionRepo.Count(db.CountOptions{
		Conditions: map[string]any{"contest_flag_id": lockedContestFlag.ID, "solved": true},
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

	solvers, currentScore, ret := CalcContestFlagState(tx, lockedContestFlag)
	if !ret.OK {
		return model.Submission{}, ret
	}
	contestFlagRepo := db.InitContestFlagRepo(tx)
	if ret = contestFlagRepo.Update(lockedContestFlag.ID, db.UpdateContestFlagOptions{
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

	score, ret := CalcTeamScore(tx, lockedTeam, contest.Blood)
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

	prometheus.RecordFlagSubmission(contest.ID, string(contestChallenge.Type), true)
	return submission, model.SuccessRetVal()
}

func CountAttempts(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) int64 {
	count, _ := db.InitSubmissionRepo(tx).Count(db.CountOptions{
		Conditions: map[string]any{"team_id": team.ID, "contest_challenge_id": contestChallenge.ID, "solved": false},
	})
	return count
}

// CheckIfSolved contestChallenge 需要预加载 ContestFlags
func CheckIfSolved(tx *gorm.DB, team model.Team, contestFlags []model.ContestFlag) bool {
	if len(contestFlags) == 0 {
		return true
	}
	solvedCount, ret := db.InitTeamFlagRepo(tx).CountSolvedForChallenge(team.ID, contestFlags[0].ContestChallengeID)
	if !ret.OK {
		return false
	}
	return solvedCount == int64(len(contestFlags))
}
