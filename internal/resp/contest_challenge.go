package resp

import (
	"CBCTF/internal/model"
	"github.com/gin-gonic/gin"
)

// GetContestChallengeResp 需要预加载 Challenge ContestFlags
func GetContestChallengeResp(contestChallenge model.ContestChallenge) gin.H {
	data := gin.H{
		"id":       contestChallenge.Challenge.RandID,
		"name":     contestChallenge.Name,
		"desc":     contestChallenge.Desc,
		"attempt":  contestChallenge.Attempt,
		"type":     contestChallenge.Type,
		"category": contestChallenge.Challenge.Category,
		"hidden":   contestChallenge.Hidden,
		"score": func() float64 {
			var score float64
			for _, flag := range contestChallenge.ContestFlags {
				score += flag.CurrentScore
			}
			return score
		}(),
		"solvers": func() int64 {
			var solvers int64
			for _, flag := range contestChallenge.ContestFlags {
				solvers += flag.Solvers
			}
			return solvers
		}(),
		"hints": contestChallenge.Hints,
		"tags":  contestChallenge.Tags,
	}
	return data
}
