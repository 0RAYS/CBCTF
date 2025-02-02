package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
)

func CreateSubmission(ctx context.Context, contestID, teamID, userID uint, challengeID, value string) (model.Submission, bool, string) {
	if _, ok, _ := GetSubmission(ctx, contestID, teamID, challengeID); ok {
		return model.Submission{}, false, "SubmissionExists"
	}
	solved := VerifyFlag(ctx, contestID, teamID, challengeID, value)
	submission := model.InitSubmission(contestID, challengeID, teamID, userID, value, solved)
	if err := DB.WithContext(ctx).Model(model.Submission{}).Create(&submission).Error; err != nil {
		return model.Submission{}, false, "CreateSubmissionError"
	}
	return submission, true, "Success"
}

func GetSubmission(ctx context.Context, contestID, teamID uint, challengeID string) (model.Submission, bool, string) {
	var submission model.Submission
	res := DB.WithContext(ctx).Model(model.Submission{}).Where("contest_id = ? AND team_id = ? AND challenge_id = ?", contestID, teamID, challengeID).Find(&submission).Limit(1)
	if res.RowsAffected != 1 {
		return model.Submission{}, false, "SubmissionNotFound"
	}
	return submission, true, "Success"
}

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
