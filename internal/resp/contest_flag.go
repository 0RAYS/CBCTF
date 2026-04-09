package resp

import (
	"CBCTF/internal/model"
	"CBCTF/internal/view"

	"github.com/gin-gonic/gin"
)

func GetContestFlagResp(contestFlag model.ContestFlag) gin.H {
	return gin.H{
		"id":                   contestFlag.ID,
		"contest_id":           contestFlag.ContestID,
		"contest_challenge_id": contestFlag.ContestChallengeID,
		"challenge_flag_id":    contestFlag.ChallengeFlagID,
		"value":                contestFlag.Value,
		"score":                contestFlag.Score,
		"current_score":        contestFlag.CurrentScore,
		"decay":                contestFlag.Decay,
		"min_score":            contestFlag.MinScore,
		"score_type":           contestFlag.ScoreType,
		"solvers":              contestFlag.Solvers,
		"last":                 contestFlag.Last,
	}
}

func GetContestFlagSolverResp(solver view.ContestFlagSolverView) gin.H {
	return gin.H{
		"user_id":   solver.UserID,
		"user_name": solver.UserName,
		"team_id":   solver.TeamID,
		"team_name": solver.TeamName,
		"score":     solver.Score,
		"solved_at": solver.SolvedAt,
	}
}
