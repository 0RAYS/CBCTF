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

// SolvedMutex 使用定时任务 cron.clearSubmissionMutex 清理锁
var SolvedMutex sync.Map

func Submit(tx *gorm.DB, user model.User, team model.Team, contest model.Contest, contestChallenge model.ContestChallenge, form dto.SubmitFlagForm, ip string) (model.Submission, model.RetVal) {
	if contestChallenge.Attempt != 0 && contestChallenge.Attempt <= CountAttempts(tx, team, contestChallenge) {
		return model.Submission{}, model.RetVal{Msg: i18n.Model.Submission.NotAllowed}
	}
	submissionRepo := db.InitSubmissionRepo(tx)
	options := db.CreateSubmissionOptions{
		ContestChallengeID: contestChallenge.ID,
		ContestID:          contest.ID,
		ChallengeID:        contestChallenge.ChallengeID,
		TeamID:             team.ID,
		UserID:             user.ID,
		Value:              form.Flag,
		Score:              team.Score,
		IP:                 ip,
	}
	solved, contestFlag, teamFlag, ret := VerifyFlag(tx, team, contestChallenge, form.Flag)
	options.ContestFlagID = contestFlag.ID
	options.Solved = solved
	if !ret.OK {
		return model.Submission{}, ret
	}
	submission, ret := submissionRepo.Create(options)
	if !ret.OK {
		return model.Submission{}, ret
	}
	if solved {
		// 正确时需要更新分数等信息, 加锁
		mu, _ := SolvedMutex.LoadOrStore(contestFlag.ID, &sync.Mutex{})
		mu.(*sync.Mutex).Lock()
		defer mu.(*sync.Mutex).Unlock()
		teamFlagRepo := db.InitTeamFlagRepo(tx)
		if ret = teamFlagRepo.Update(teamFlag.ID, db.UpdateTeamFlagRepo{Solved: &solved}); !ret.OK {
			return model.Submission{}, ret
		}
		_, currentScore, ret := CalcContestFlagState(tx, contestFlag)
		if !ret.OK {
			return model.Submission{}, ret
		}
		contestFlagRepo := db.InitContestFlagRepo(tx)
		if ret = contestFlagRepo.DiffUpdate(contestFlag.ID, db.DiffUpdateContestFlagOptions{Solvers: 1}); !ret.OK {
			return model.Submission{}, ret
		}
		if ret = contestFlagRepo.Update(contestFlag.ID, db.UpdateContestFlagOptions{
			CurrentScore: &currentScore,
			Last:         &submission.CreatedAt,
		}); !ret.OK {
			return model.Submission{}, ret
		}
		score, ret := CalcTeamScore(tx, team, contest.Blood)
		if !ret.OK {
			return model.Submission{}, ret
		}
		teamRepo := db.InitTeamRepo(tx)
		if ret = teamRepo.Update(team.ID, db.UpdateTeamOptions{
			Score: &score,
			Last:  &submission.CreatedAt,
		}); !ret.OK {
			return model.Submission{}, ret
		}
		if ret = submissionRepo.Update(submission.ID, db.UpdateSubmissionOptions{Score: &score}); !ret.OK {
			return model.Submission{}, ret
		}
	}
	prometheus.UpdateFlagSubmissionMetrics(contest, contestChallenge, team, solved)
	return submission, model.SuccessRetVal()
}

func CountAttempts(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) int64 {
	submissionRepo := db.InitSubmissionRepo(tx)
	count, _ := submissionRepo.Count(db.CountOptions{
		Conditions: map[string]any{"team_id": team.ID, "contest_challenge_id": contestChallenge.ID, "solved": false},
	})
	return count
}

// CheckIfSolved contestChallenge 需要预加载 ContestFlags
func CheckIfSolved(tx *gorm.DB, team model.Team, contestFlags []model.ContestFlag) bool {
	if len(contestFlags) == 0 {
		return true
	}
	submissionRepo := db.InitSubmissionRepo(tx)
	count, _ := submissionRepo.Count(db.CountOptions{
		Conditions: map[string]any{"team_id": team.ID, "contest_challenge_id": contestFlags[0].ContestChallengeID, "solved": true},
	})
	return count == int64(len(contestFlags))
}
