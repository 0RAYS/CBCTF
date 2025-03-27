package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

// GetUsageResp model.Usage 需要预加载
func GetUsageResp(usage model.Usage) gin.H {
	return gin.H{
		"id":       usage.Challenge.ID,
		"name":     usage.Challenge.Name,
		"category": usage.Challenge.Category,
		"type":     usage.Challenge.Type,
		"score": func() float64 {
			var score float64
			for _, flag := range usage.Flags {
				score += flag.CurrentScore
			}
			return score
		}(),
		"solvers": func() int64 {
			var solvers int64
			for _, flag := range usage.Flags {
				solvers += flag.Solvers
			}
			return solvers
		},
		"hints": usage.Hints,
		"tags":  usage.Tags,
	}
}
