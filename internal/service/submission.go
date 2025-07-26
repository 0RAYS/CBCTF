package service

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	db "CBCTF/internal/repo"
	"gorm.io/gorm"
	"sync"
)

// SolvedMutex 使用定时任务 cron.ClearContestChallengeMutex 清理锁
var SolvedMutex sync.Map

// Submit model.Usage 需要预加载
func Submit(tx *gorm.DB, user model.User, team model.Team, contestChallenge model.ContestChallenge, form f.SubmitFlagForm, ip string) (string, model.Submission, bool, string) {
	if contestChallenge.Attempt != 0 && contestChallenge.Attempt <= CountAttempts(tx, team, contestChallenge) {
		return "", model.Submission{}, false, i18n.NotAllowSubmit
	}
	submissionRepo := db.InitSubmissionRepo(tx)
	options := db.CreateSubmissionOptions{
		ContestChallengeID: contestChallenge.ID,
		ContestID:          team.ContestID,
		ChallengeID:        contestChallenge.ChallengeID,
		TeamID:             team.ID,
		UserID:             user.ID,
		Value:              form.Flag,
		Score:              team.Score,
		IP:                 ip,
	}
	solved, contestFlag, teamFlag, ok, result := VerifyFlag(tx, team, contestChallenge, form.Flag)
	options.ContestFlagID = contestFlag.ID
	options.Solved = solved
	if !ok {
		return "", model.Submission{}, false, result
	}
	submission, ok, msg := submissionRepo.Create(options)
	if !ok {
		return "", model.Submission{}, false, msg
	}
	if solved {
		teamFlagRepo := db.InitTeamFlagRepo(tx)
		if ok, msg = teamFlagRepo.Update(teamFlag.ID, db.UpdateTeamFlagRepo{Solved: &solved}); !ok {
			return "", model.Submission{}, false, msg
		}
		// 正确时需要更新分数等信息, 加锁
		mu, _ := SolvedMutex.LoadOrStore(contestFlag.ID, &sync.Mutex{})
		mu.(*sync.Mutex).Lock()
		defer mu.(*sync.Mutex).Unlock()

		solvers, currentScore, ok, msg := CalcContestFlagState(tx, contestFlag)
		if !ok {
			return "", model.Submission{}, false, msg
		}
		contestFlagRepo := db.InitContestFlagRepo(tx)
		if ok, msg = contestFlagRepo.Update(contestFlag.ID, db.UpdateContestFlagOptions{
			CurrentScore: &currentScore,
			Solvers:      &solvers,
			Last:         &submission.CreatedAt,
		}); !ok {
			return "", model.Submission{}, false, msg
		}
		score, ok, msg := CalcTeamScore(tx, team)
		if !ok {
			return "", model.Submission{}, false, msg
		}
		teamRepo := db.InitTeamRepo(tx)
		if ok, msg = teamRepo.Update(team.ID, db.UpdateTeamOptions{
			Score: &score,
			Last:  &submission.CreatedAt,
		}); !ok {
			return "", model.Submission{}, false, msg
		}
		if ok, msg := submissionRepo.Update(submission.ID, db.UpdateSubmissionOptions{Score: &score}); !ok {
			return "", model.Submission{}, false, msg
		}
	}
	return result, submission, true, i18n.Success
}

func CountAttempts(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) int64 {
	submissionRepo := db.InitSubmissionRepo(tx)
	count, _, _ := submissionRepo.Count(db.CountOptions{
		Conditions: map[string]any{
			"team_id":              team.ID,
			"contest_challenge_id": contestChallenge.ID,
			"solved":               false,
		},
	})
	return count
}

// CheckIfSolved contestChallenge 需要预加载 ContestFlags
func CheckIfSolved(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) bool {
	submissionRepo := db.InitSubmissionRepo(tx)
	count, _, _ := submissionRepo.Count(db.CountOptions{
		Conditions: map[string]any{
			"team_id":              team.ID,
			"contest_challenge_id": contestChallenge.ID,
			"solved":               true,
		},
	})
	return count == int64(len(contestChallenge.ContestFlags))
}
