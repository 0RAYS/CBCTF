package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"gorm.io/gorm"
)

// CreateSubmission is a function to create a new submission
func CreateSubmission(tx *gorm.DB, ctx context.Context, contest model.Contest, team model.Team, user model.User, challenge model.Challenge, value string) (model.Submission, bool, string) {
	if IsSolved(ctx, contest, team, challenge) {
		return model.Submission{}, false, "AlreadySolved"
	}
	usage, ok, _ := GetUsageBy2ID(ctx, contest.ID, challenge.ID)
	if !ok || usage.Attempt <= CountAttempts(ctx, contest, team, challenge) || !contest.IsRunning() {
		return model.Submission{}, false, "NotAllowSubmit"
	}
	solved := VerifyFlag(ctx, contest.ID, team.ID, challenge.ID, value)
	if solved {
		if ok, msg := AddSolvers(tx, usage.ID); !ok {
			return model.Submission{}, false, msg
		}
	}
	submission := model.InitSubmission(contest.ID, challenge.ID, team.ID, user.ID, value, solved)
	if err := tx.Model(model.Submission{}).Create(&submission).Error; err != nil {

		return model.Submission{}, false, "CreateSubmissionError"
	}
	return submission, true, "Success"
}

// GetTeamSubmissions is a function to get submission
func GetTeamSubmissions(ctx context.Context, contest model.Contest, team model.Team, challenge model.Challenge) ([]model.Submission, bool, string) {
	var submissions []model.Submission
	res := DB.WithContext(ctx).Model(model.Submission{}).Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contest.ID, team.ID, challenge.ID).Find(&submissions)
	if res.RowsAffected != 1 {
		return []model.Submission{}, false, "SubmissionNotFound"
	}
	return submissions, true, "Success"
}

func IsSolved(ctx context.Context, contest model.Contest, team model.Team, challenge model.Challenge) bool {
	var submission model.Submission
	res := DB.WithContext(ctx).Model(model.Submission{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ? AND solved = ?", contest.ID, team.ID, challenge.ID, true).Find(&submission)
	if res.RowsAffected != 1 {
		return false
	}
	return true
}

func CountAttempts(ctx context.Context, contest model.Contest, team model.Team, challenge model.Challenge) int64 {
	var count int64
	res := DB.WithContext(ctx).Model(model.Submission{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contest.ID, team.ID, challenge.ID).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count attempts: %v", res.Error)
		return 0
	}
	return count
}

// GetSubmissions is a function to get submissions
func GetSubmissions(ctx context.Context, limit, offset int) ([]model.Submission, int64, bool, string) {
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}
	var submissions []model.Submission
	var count int64
	res := DB.WithContext(ctx).Model(model.Submission{})
	if res.Count(&count).Error != nil {
		log.Logger.Warningf("Failed to count submissions: %v", res.Error)
		return nil, 0, false, "UnknownError"
	}
	if res = res.Order("created_at desc").Limit(limit).Offset(offset).Find(&submissions); res.Error != nil {
		log.Logger.Warningf("Failed to get submissions: %v", res.Error)
		return nil, 0, false, "UnknownError"
	}
	return submissions, count, true, "Success"
}

func DeleteSubmission(tx *gorm.DB, id uint) (bool, string) {
	var submission model.Submission
	res := tx.Model(model.Submission{}).Where("id = ?", id).Delete(&submission)
	if res.RowsAffected != 1 {
		return false, "SubmissionNotFound"
	}
	return true, "Success"
}
