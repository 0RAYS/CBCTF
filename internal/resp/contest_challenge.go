package resp

import (
	"CBCTF/internal/model"
	"CBCTF/internal/view"

	"github.com/gin-gonic/gin"
)

func GetVictimStatusResp(status view.VictimStatusView) gin.H {
	return gin.H{
		"target":    status.Targets,
		"duration":  status.Duration,
		"remaining": status.Remaining,
		"status":    status.Status,
	}
}

func getContestChallengeBaseResp(contestChallenge model.ContestChallenge) gin.H {
	data := gin.H{
		"id":          contestChallenge.Challenge.RandID,
		"name":        contestChallenge.Name,
		"description": contestChallenge.Description,
		"attempt":     contestChallenge.Attempt,
		"type":        contestChallenge.Type,
		"category":    contestChallenge.Category,
		"hidden":      contestChallenge.Hidden,
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

func GetContestChallengeResp(contestChallengeView view.ContestChallengeView) gin.H {
	data := getContestChallengeBaseResp(contestChallengeView.ContestChallenge)
	data["attempts"] = contestChallengeView.Attempts
	data["init"] = contestChallengeView.Init
	data["solved"] = contestChallengeView.Solved
	data["remote"] = GetVictimStatusResp(contestChallengeView.Remote)
	data["file"] = contestChallengeView.FileName
	return data
}

func GetAdminContestChallengeResp(contestChallenge model.ContestChallenge) gin.H {
	return getContestChallengeBaseResp(contestChallenge)
}

func GetContestChallengeStatusResp(status view.ContestChallengeStatusView) gin.H {
	return gin.H{
		"attempts": status.Attempts,
		"init":     status.Init,
		"solved":   status.Solved,
		"remote":   GetVictimStatusResp(status.Remote),
		"file":     status.FileName,
	}
}
