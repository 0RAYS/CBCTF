package resp

import (
	"CBCTF/internal/view"

	"github.com/gin-gonic/gin"
)

func GetSolvedStateViewResp(states []view.ScoreboardSolvedStateView) []gin.H {
	data := make([]gin.H, 0, len(states))
	for _, state := range states {
		data = append(data, gin.H{
			"category": state.Category,
			"solved":   state.Solved,
			"all":      state.All,
		})
	}
	return data
}

func GetTeamRankingResp(teamView view.TeamRankingView, isAdmin bool) gin.H {
	team := teamView.Team
	data := gin.H{
		"id":          team.ID,
		"name":        team.Name,
		"description": team.Description,
		"score":       team.Score,
		"picture":     team.Picture,
		"last":        team.Last,
		"hidden":      team.Hidden,
		"captain_id":  team.CaptainID,
		"solved":      GetSolvedStateViewResp(teamView.Solved),
	}
	if isAdmin {
		data["captcha"] = team.Captcha
	}
	return data
}

func GetScoreboardResp(teams []view.ScoreboardTeamView) []gin.H {
	data := make([]gin.H, 0)
	for _, teamView := range teams {
		team := teamView.Team
		solved := make([]gin.H, 0)
		for _, challenge := range teamView.Challenges {
			solved = append(solved, gin.H{
				"id":       challenge.ID,
				"total":    challenge.Total,
				"solved":   challenge.Solved,
				"name":     challenge.Name,
				"category": challenge.Category,
			})
		}
		data = append(data, gin.H{
			"id":          team.ID,
			"name":        team.Name,
			"description": team.Description,
			"score":       team.Score,
			"picture":     team.Picture,
			"last":        team.Last,
			"users":       teamView.UserCount,
			"challenges":  solved,
		})
	}
	return data
}

func GetRankTimelineResp(teams []view.RankTimelineTeamView) []gin.H {
	data := make([]gin.H, 0)
	for _, team := range teams {
		timeline := make([]gin.H, 0)
		for _, point := range team.Timeline {
			timeline = append(timeline, gin.H{
				"time":  point.Time,
				"score": point.Score,
			})
		}
		data = append(data, gin.H{
			"id":       team.ID,
			"name":     team.Name,
			"picture":  team.Picture,
			"rank":     team.Rank,
			"score":    team.Score,
			"timeline": timeline,
		})
	}
	return data
}
