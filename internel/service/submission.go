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
func Submit(tx *gorm.DB, contest model.Contest, user model.User, team model.Team, usage model.Usage, form f.SubmitFlagForm) (model.Submission, bool, string) {
	if usage.Attempt != 0 && usage.Attempt <= CountAttempts(tx, team, usage) {
		return model.Submission{}, false, "NotAllowSubmit"
	}
	submissionRepo := db.InitSubmissionRepo(tx)
	options := db.CreateSubmissionOptions{
		UsageID:     usage.ID,
		ContestID:   team.ContestID,
		ChallengeID: usage.ChallengeID,
		TeamID:      team.ID,
		UserID:      user.ID,
		Value:       form.Flag,
		Solved:      false,
		Score:       team.Score,
	}
	solved, flag, ok, msg := VerifyFlag(tx, team, usage, form.Flag)
	options.FlagID = flag.ID
	if !ok {
		if options.FlagID > 0 {
			submissionRepo.Create(options)
		}
		return model.Submission{}, false, msg
	}
	submission, ok, msg := submissionRepo.Create(options)
	if solved {
		submission.Solved = true
		submissionRepo.Update(submission.ID, db.UpdateSubmissionOptions{Solved: &solved})
		answerRepo := db.InitAnswerRepo(tx)
		answerRepo.Update(flag.ID, db.UpdateAnswerOptions{Solved: &solved})
		// 正确时需要更新分数等信息, 加锁
		mu, _ := SolvedMutex.LoadOrStore(usage.ID, &sync.Mutex{})
		mu.(*sync.Mutex).Lock()
		defer mu.(*sync.Mutex).Unlock()

		solvers, currentScore, ok, msg := CalcSolversAndScore(tx, flag)
		if !ok {
			return model.Submission{}, false, msg
		}
		rate, blood := flag.CalcBlood(team.ID)
		if !contest.Blood {
			rate = 0
		}
		if blood >= 0 && blood <= 2 {
			flag.Blood[blood] = team.ID
		}
		score, ok, msg := CalcTeamScore(tx, team.ID)
		if !ok {
			return model.Submission{}, false, msg
		}
		score += currentScore + flag.Score*rate
		teamRepo := db.InitTeamRepo(tx)
		ok, msg = teamRepo.Update(team.ID, db.UpdateTeamOptions{
			Score: &score,
			Last:  &submission.CreatedAt,
		})
		if !ok {
			return model.Submission{}, false, msg
		}
		flagRepo := db.InitFlagRepo(tx)
		ok, msg = flagRepo.Update(flag.ID, db.UpdateFlagOptions{
			CurrentScore: &currentScore,
			Solvers:      &solvers,
			Blood:        &flag.Blood,
		})
		if !ok {
			return model.Submission{}, false, msg
		}
	}
	return submission, true, "Success"
}

// IsSolved model.Usage 需要预加载
func IsSolved(tx *gorm.DB, team model.Team, usage model.Usage) bool {
	var (
		count                 int
		submissionRepo        = db.InitSubmissionRepo(tx)
		submissions, _, ok, _ = submissionRepo.GetAllByKeyID("team_id", team.ID, -1, -1, true, false)
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
	submissions, _, ok, _ := submissionRepo.GetAllByKeyID("team_id", team.ID, -1, -1, false, false)
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

// CountFlagSolved 统计指定 model.Flag 的解题次数
func CountFlagSolved(tx *gorm.DB, flag model.Flag) (int64, bool, string) {
	var (
		count                   int64
		submissionRepo          = db.InitSubmissionRepo(tx)
		submissions, _, ok, msg = submissionRepo.GetAllByKeyID("contest_id", flag.ContestID, -1, -1, true, true)
	)
	if !ok {
		return count, false, msg
	}
	for _, submission := range submissions {
		if submission.FlagID == flag.ID {
			count++
		}
	}
	if count < flag.Solvers {
		// 不考虑更新失败的情况, 不回滚
		flagRepo := db.InitFlagRepo(tx)
		flagRepo.Update(flag.ID, db.UpdateFlagOptions{Solvers: &count})
	}
	return count, true, "Success"
}
