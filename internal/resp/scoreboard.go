package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetTeamRankingResp(teamsData []struct {
	Team   model.Team
	Solved []model.ContestFlag
}, flags []model.ContestFlag, admin bool) gin.H {
	data := make([]gin.H, 0)
	for _, team := range teamsData {
		tmp := gin.H{
			"id":     team.Team.ID,
			"name":   team.Team.Name,
			"desc":   team.Team.Desc,
			"score":  team.Team.Score,
			"avatar": team.Team.Avatar,
			"last":   team.Team.Last,
			"users":  len(team.Team.Users),
			"solved": GetSolvedStateResp(team.Solved, flags),
		}
		if admin {
			tmp["hidden"] = team.Team.Hidden
		}
		data = append(data, tmp)
	}
	return gin.H{"teams": data}
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
			"id":         team.ID,
			"name":       team.Name,
			"desc":       team.Desc,
			"score":      team.Score,
			"avatar":     team.Avatar,
			"last":       team.Last,
			"users":      len(team.Users),
			"challenges": solved,
		})
	}
	return data
}

// GetRankTimelineResp 需要预加载 model.Submission
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
			"avatar":   team.Avatar,
			"rank":     team.Rank,
			"score":    team.Score,
			"timeline": timeline,
		})
	}
	return data
}
