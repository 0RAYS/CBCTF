package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
	"CBCTF/internal/view"
)

func GetChallengeFlagResp(flag model.ChallengeFlag) gin.H {
	return gin.H{
		"id":            flag.ID,
		"challenge_id":  flag.ChallengeID,
		"value":         flag.Value,
		"binding":       flag.Binding,
		"generation":    flag.Generation,
		"active":        flag.Active,
		"score":         flag.Score,
		"current_score": flag.CurrentScore,
		"decay":         flag.Decay,
		"min_score":     flag.MinScore,
		"score_type":    flag.ScoreType,
		"solvers":       flag.Solvers,
		"last":          flag.Last,
	}
}

func GetFlagSolverResp(solver view.FlagSolverView) gin.H {
	return gin.H{
		"user_id":   solver.UserID,
		"user_name": solver.UserName,
		"team_id":   solver.TeamID,
		"team_name": solver.TeamName,
		"score":     solver.Score,
		"solved_at": solver.SolvedAt,
	}
}
