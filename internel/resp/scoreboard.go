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

func GetScoreboardResp(globalMap map[uint]int, teamMap map[uint]map[uint]int, teams []model.Team) gin.H {
	data := make([]gin.H, 0)
	for _, team := range teams {
		solved := make([]gin.H, 0)
		for challengeID, count := range teamMap[team.ID] {
			status := 0
			if count == 0 {
				status = 0 // 未解
			} else if count > 0 && count < globalMap[challengeID] {
				status = 1 // 部分解
			} else {
				status = 2 // 完全解
			}
			solved = append(solved, gin.H{
				"id":     challengeID,
				"status": status,
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
	return gin.H{"flags": data}
}
