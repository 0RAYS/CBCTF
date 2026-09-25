package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

func GetGeneratorResp(generator model.Generator) gin.H {
	return gin.H{
		"id":             generator.ID,
		"challenge_id":   generator.ChallengeID,
		"challenge_name": generator.ChallengeName,
		"contest_id":     generator.ContestID.V,
		"name":           generator.Name,
		"image":          generator.Image,
		"start_time":     generator.CreatedAt,
		"success":        generator.Success,
		"success_last":   generator.SuccessLast,
		"failure":        generator.Failure,
		"failure_last":   generator.FailureLast,
		"status":         generator.Status,
	}
}
