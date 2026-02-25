package resp

import (
	"CBCTF/internal/model"
	"crypto/rand"
	"math/big"

	"github.com/gin-gonic/gin"
)

// GetContestChallengeResp model.ContestChallenge Preload model.Challenge model.ContestFlag
func GetContestChallengeResp(contestChallenge model.ContestChallenge) gin.H {
	options := make([]gin.H, 0)
	for _, option := range contestChallenge.Challenge.Options {
		options = append(options, gin.H{
			"rand_id": option.RandID,
			"content": option.Content,
		})
	}
	for i := len(options) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			panic(err)
		}
		options[i], options[j.Int64()] = options[j.Int64()], options[i]
	}
	data := gin.H{
		"id":          contestChallenge.Challenge.RandID,
		"name":        contestChallenge.Name,
		"description": contestChallenge.Description,
		"attempt":     contestChallenge.Attempt,
		"type":        contestChallenge.Type,
		"category":    contestChallenge.Category,
		"hidden":      contestChallenge.Hidden,
		"options":     options,
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
