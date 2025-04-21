package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetFlagResp(flag model.Flag) gin.H {
	return gin.H{
		"id":            flag.ID,
		"value":         flag.Value,
		"score":         flag.Score,
		"current_score": flag.CurrentScore,
		"decay":         flag.Decay,
		"min_score":     flag.MinScore,
		"score_type":    flag.ScoreType,
		"solvers":       flag.Solvers,
		"last":          flag.Last,
	}
}
