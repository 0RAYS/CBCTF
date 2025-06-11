package resp

import (
	"CBCTF/internel/model"
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
			status := 0
			if count == 0 {
				status = 0 // 未解
			} else if count > 0 && count < globalMap[challengeRandID] {
				status = 1 // 部分解
			} else {
				status = 2 // 完全解
			}
			solved = append(solved, gin.H{
				"id":       challengeRandID,
				"status":   status,
				"name":     challengeMap[challengeRandID].Name,
				"category": challengeMap[challengeRandID].Category,
			})
		}
		data = append(data, gin.H{
			"id":     team.ID,
			"name":   team.Name,
			"score":  team.Score,
			"avatar": team.Avatar,
			"solved": solved,
		})
	}
	return data
}
