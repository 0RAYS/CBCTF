package resp

import (
	"CBCTF/internal/model"
	"math/rand"

	"github.com/gin-gonic/gin"
)

// GetContestChallengeResp 需要预加载 Challenge ContestFlags
func GetContestChallengeResp(contestChallenge model.ContestChallenge) gin.H {
	options := make([]gin.H, 0)
	for _, option := range contestChallenge.Challenge.Options {
		options = append(options, gin.H{
			"rand_id": option.RandID,
			"content": option.Content,
		})
	}
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})
	data := gin.H{
		"id":       contestChallenge.Challenge.RandID,
		"name":     contestChallenge.Name,
		"desc":     contestChallenge.Desc,
		"attempt":  contestChallenge.Attempt,
		"type":     contestChallenge.Type,
		"category": contestChallenge.Challenge.Category,
		"hidden":   contestChallenge.Hidden,
		"options":  options,
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
