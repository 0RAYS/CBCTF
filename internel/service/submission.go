package service

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

func CountAttempts(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) int64 {
	submissionRepo := db.InitSubmissionRepo(tx)
	count, _, _ := submissionRepo.CountWithConditions(db.GetOptions{
		{Key: "team_id", Value: team.ID, Op: "and"},
		{Key: "contest_challenge_id", Value: contestChallenge.ID, Op: "and"},
		{Key: "solved", Value: false, Op: "and"},
	})
	return count
}

// CheckIfSolved contestChallenge 需要预加载 ContestFlags
func CheckIfSolved(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) bool {
	submissionRepo := db.InitSubmissionRepo(tx)
	count, _, _ := submissionRepo.CountWithConditions(db.GetOptions{
		{Key: "team_id", Value: team.ID, Op: "and"},
		{Key: "contest_challenge_id", Value: contestChallenge.ID, Op: "and"},
		{Key: "solved", Value: true, Op: "and"},
	})
	return count == int64(len(contestChallenge.ContestFlags))
}
