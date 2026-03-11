package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetGeneratorResp(generator model.Generator) gin.H {
	return gin.H{
		"id":           generator.ID,
		"challenge_id": generator.ChallengeID,
		"contest_id":   generator.ContestID,
		"name":         generator.Name,
		"start_time":   generator.CreatedAt,
		"success":      generator.Success,
		"success_last": generator.SuccessLast,
		"failure":      generator.Failure,
		"failure_last": generator.FailureLast,
	}
}
