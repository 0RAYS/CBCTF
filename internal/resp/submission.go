package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetSubmissionResp(submission model.Submission) gin.H {
	return gin.H{
		"id":                   submission.ID,
		"contest_id":           submission.ContestID,
		"contest_challenge_id": submission.ContestChallengeID,
		"challenge_id":         submission.ChallengeID,
		"team_id":              submission.TeamID,
		"user_id":              submission.UserID,
		"solved":               submission.Solved,
		"value":                submission.Value,
		"score":                submission.Score,
		"ip":                   submission.IP,
		"created_at":           submission.CreatedAt,
	}
}
