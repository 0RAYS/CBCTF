package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
)

// CreateSubmission is a function to create a new submission
func CreateSubmission(ctx context.Context, contestID, teamID, userID uint, challengeID, value string) (model.Submission, bool, string) {
	if IsSolved(ctx, contestID, teamID, challengeID) {
		return model.Submission{}, false, "AlreadySolved"
	}
	solved := VerifyFlag(ctx, contestID, teamID, challengeID, value)
	submission := model.InitSubmission(contestID, challengeID, teamID, userID, value, solved)
	if err := DB.WithContext(ctx).Model(model.Submission{}).Create(&submission).Error; err != nil {
		return model.Submission{}, false, "CreateSubmissionError"
	}
	return submission, true, "Success"
}

// GetTeamSubmissions is a function to get submission
func GetTeamSubmissions(ctx context.Context, contestID, teamID uint, challengeID string) ([]model.Submission, bool, string) {
	var submissions []model.Submission
	res := DB.WithContext(ctx).Model(model.Submission{}).Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contestID, teamID, challengeID).Find(&submissions)
	if res.RowsAffected != 1 {
		return []model.Submission{}, false, "SubmissionNotFound"
	}
	return submissions, true, "Success"
}

func IsSolved(ctx context.Context, contestID, teamID uint, challengeID string) bool {
	var submission model.Submission
	res := DB.WithContext(ctx).Model(model.Submission{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ? AND solved = ?", contestID, teamID, challengeID, true).Find(&submission)
	if res.RowsAffected != 1 {
		return false
	}
	return true
}

func CountAttempts(ctx context.Context, contestID, teamID uint, challengeID string) int64 {
	var count int64
	res := DB.WithContext(ctx).Model(model.Submission{}).
		Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contestID, teamID, challengeID).Count(&count)
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
