package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
	"CBCTF/internal/view"
)

// GetSolvedStateResp 按题目分类统计解题数量，两组 flag 都需预加载 Challenge
func GetSolvedStateResp(solved []model.ChallengeFlag, all []model.ChallengeFlag) []gin.H {
	return GetSolvedStateViewResp(view.BuildSolvedStateViews(solved, all))
}

func GetTeamResp(teamView view.TeamView, isAdmin bool) gin.H {
	team := teamView.Team
	data := gin.H{
		"id":          team.ID,
		"contest_id":  team.ContestID,
		"name":        team.Name,
		"score":       team.Score,
		"picture":     team.Picture,
		"last":        team.Last,
		"rank":        team.Rank,
		"users":       teamView.UserCount,
		"description": team.Description,
		"captain_id":  team.CaptainID,
		"banned":      team.Banned,
		"hidden":      team.Hidden,
	}
	if isAdmin {
		data["captcha"] = team.Captcha
		data["adjustment"] = teamView.Adjustment
	}
	markUnavailable(data, teamView.Unavailable)
	return data
}
