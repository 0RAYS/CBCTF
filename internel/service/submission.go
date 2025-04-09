package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
	"sync"
)

// SolvedMutex 使用定时任务 cron.ClearUsageMutex 清理锁
var SolvedMutex sync.Map

// Submit model.Usage 需要预加载
func Submit(tx *gorm.DB, user model.User, team model.Team, usage model.Usage, form f.SubmitFlagForm) (string, model.Submission, bool, string) {
	if usage.Attempt != 0 && usage.Attempt <= CountAttempts(tx, team, usage) {
		return "", model.Submission{}, false, "NotAllowSubmit"
	}
	submissionRepo := db.InitSubmissionRepo(tx)
	options := db.CreateSubmissionOptions{
		UsageID:     usage.ID,
		ContestID:   team.ContestID,
		ChallengeID: usage.ChallengeID,
		TeamID:      team.ID,
		UserID:      user.ID,
		Value:       form.Flag,
		Score:       team.Score,
	}
	solved, flag, answer, ok, result := VerifyFlag(tx, team, usage, form.Flag)
	options.FlagID = flag.ID
	options.Solved = solved
	if !ok {
		return "", model.Submission{}, false, result
	}
	submission, ok, msg := submissionRepo.Create(options)
	if !ok {
		return "", model.Submission{}, false, msg
	}
	if solved {
		answerRepo := db.InitAnswerRepo(tx)
		if ok, msg := answerRepo.Update(answer.ID, db.UpdateAnswerOptions{Solved: &solved}); !ok {
			return "", model.Submission{}, false, msg
		}
		// 正确时需要更新分数等信息, 加锁
		mu, _ := SolvedMutex.LoadOrStore(flag.ID, &sync.Mutex{})
		mu.(*sync.Mutex).Lock()
		defer mu.(*sync.Mutex).Unlock()

		solvers, currentScore, ok, msg := CalcSolversAndScore(tx, flag)
		if !ok {
			return "", model.Submission{}, false, msg
		}
		_, blood := flag.CalcBlood(team.ID)
		if blood >= 0 && blood <= 2 {
			flag.Blood[blood] = team.ID
		}
		flagRepo := db.InitFlagRepo(tx)
		if ok, msg = flagRepo.Update(flag.ID, db.UpdateFlagOptions{
			CurrentScore: &currentScore,
			Solvers:      &solvers,
			Blood:        &flag.Blood,
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
	return result, submission, true, "Success"
}

// IsSolved model.Usage 需要预加载
func IsSolved(tx *gorm.DB, team model.Team, usage model.Usage) bool {
	var (
		count                 int
		submissionRepo        = db.InitSubmissionRepo(tx)
		submissions, _, ok, _ = submissionRepo.GetAllByKeyID("team_id", team.ID, -1, -1, true)
	)
	if !ok {
		return false
	}
	for _, submission := range submissions {
		if submission.UsageID == usage.ID {
			count++
		}
	}
	if count != len(usage.Flags) {
		return false
	}
	return true
}

// CountAttempts 统计题目的尝试次数
func CountAttempts(tx *gorm.DB, team model.Team, usage model.Usage) int64 {
	var count int64
	submissionRepo := db.InitSubmissionRepo(tx)
	submissions, _, ok, _ := submissionRepo.GetAllByKeyID("team_id", team.ID, -1, -1, false)
	if !ok {
		return count
	}
	for _, submission := range submissions {
		if submission.UsageID == usage.ID {
			count++
		}
	}
	return count
}
