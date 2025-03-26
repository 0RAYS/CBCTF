package resp

import (
	"CBCTF/internel/config"
	"CBCTF/internel/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetTeamRankingResp(teamsData []struct {
	Team   model.Team
	Solved []model.Flag
}, flags []model.Flag) gin.H {
	categories := make(map[uint]string)
	for _, v := range flags {
		categories[v.UsageID] = v.Usage.Challenge.Category
	}
	allCount := make(map[string]int64)
	for _, v := range flags {
		allCount[v.Usage.Challenge.Category] += 1
	}
	data := make([]gin.H, 0)
	for _, team := range teamsData {
		solvedCount := make(map[string]int64)
		for _, flag := range team.Solved {
			solvedCount[categories[flag.UsageID]] += 1
		}
		for k, _ := range allCount {
			if _, ok := solvedCount[k]; !ok {
				solvedCount[k] = 0
			}
		}
		data = append(data, gin.H{
			"name":   team.Team.Name,
			"score":  team.Team.Score,
			"avatar": fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(team.Team.Avatar, "/")),
			"last":   team.Team.Last,
			"solved": func() []gin.H {
				var solved []gin.H
				for k, v := range solvedCount {
					solved = append(solved, gin.H{"category": k, "solved": v, "all": allCount[k]})
				}
				return solved
			}(),
		})
	}
	return gin.H{"teams": data}
}
