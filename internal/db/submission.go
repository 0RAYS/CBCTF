package db

import (
	"CBCTF/internal/model"
	"context"
)

func CreateSubmission(ctx context.Context, contestID, teamID, userID uint, challengeID, value string) (model.Submission, bool, string) {
	solved := VerifyFlag(ctx, contestID, teamID, challengeID, value)
	submission := model.InitSubmission(contestID, challengeID, teamID, userID, value, solved)
	if err := DB.WithContext(ctx).Model(model.Submission{}).Create(&submission).Error; err != nil {
		return model.Submission{}, false, "CreateSubmissionError"
	}
	return submission, true, "Success"
}
