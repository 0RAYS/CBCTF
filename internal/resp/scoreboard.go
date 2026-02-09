package resp

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

// GetTeamRankingResp model.ContestFlag Preload model.ContestChallenge
func GetTeamRankingResp(team model.Team, solved []model.ContestFlag, flags []model.ContestFlag, isAdmin bool) gin.H {
	data := gin.H{
		"id":          team.ID,
		"name":        team.Name,
		"description": team.Description,
		"score":       team.Score,
		"picture":     team.Picture,
		"last":        team.Last,
		"hidden":      team.Hidden,
		"captain_id":  team.CaptainID,
		"solved":      GetSolvedStateResp(solved, flags),
	}
	if isAdmin {
		data["captcha"] = team.Captcha
	}
	return data
}

func GetScoreboardResp(challengeMap map[string]model.Challenge, globalMap map[string]int, teamMap map[uint]map[string]int, teams []model.Team) []gin.H {
	data := make([]gin.H, 0)
	for _, team := range teams {
		solved := make([]gin.H, 0)
		for challengeRandID, count := range teamMap[team.ID] {
			solved = append(solved, gin.H{
				"id":       challengeRandID,
				"total":    globalMap[challengeRandID],
				"solved":   count,
				"name":     challengeMap[challengeRandID].Name,
				"category": challengeMap[challengeRandID].Category,
			})
		}
		data = append(data, gin.H{
			"id":          team.ID,
			"name":        team.Name,
			"description": team.Description,
			"score":       team.Score,
			"picture":     team.Picture,
			"last":        team.Last,
			"users":       db.InitTeamRepo(db.DB).CountAssociation(team, "Users"),
			"challenges":  solved,
		})
	}
	return data
}

// GetRankTimelineResp model.Team Preload model.Submission
func GetRankTimelineResp(teams []model.Team) []gin.H {
	data := make([]gin.H, 0)
	for _, team := range teams {
		timeline := make([]gin.H, 0)
		for _, submission := range team.Submissions {
			timeline = append(timeline, gin.H{
				"time":  submission.CreatedAt,
				"score": submission.Score,
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
