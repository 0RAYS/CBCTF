package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
	"CBCTF/internal/view"
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
	var score float64
	var solvers int64
	for _, flag := range contestChallenge.ContestFlags {
		score += flag.CurrentScore
		solvers += flag.Solvers
	}
	return gin.H{
		"id":          contestChallenge.Challenge.RandID,
		"name":        contestChallenge.Name,
		"description": contestChallenge.Description,
		"attempt":     contestChallenge.Attempt,
		"type":        contestChallenge.Type,
		"category":    contestChallenge.Category,
		"hidden":      contestChallenge.Hidden,
		"score":       score,
		"solvers":     solvers,
		"hints":       contestChallenge.Hints,
		"tags":        contestChallenge.Tags,
	}
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
