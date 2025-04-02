package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetSubmissionResp(submission model.Submission) gin.H {
	return gin.H{
		"challenge_id": submission.ChallengeID,
		"team_id":      submission.TeamID,
		"user_id":      submission.UserID,
		"solved":       submission.Solved,
		"value":        submission.Value,
	}
}
