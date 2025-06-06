package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetContestFlagResp(contestFlag model.ContestFlag) gin.H {
	return gin.H{
		"id":            contestFlag.ID,
		"value":         contestFlag.Value,
		"score":         contestFlag.Score,
		"current_score": contestFlag.CurrentScore,
		"decay":         contestFlag.Decay,
		"min_score":     contestFlag.MinScore,
		"score_type":    contestFlag.ScoreType,
		"solvers":       contestFlag.Solvers,
		"last":          contestFlag.Last,
	}
}
