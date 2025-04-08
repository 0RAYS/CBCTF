package resp

import (
	"CBCTF/internel/config"
	"CBCTF/internel/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetSolvedStateResp(solved []model.Flag, flags []model.Flag) []gin.H {
	categories := make(map[uint]string)
	for _, v := range flags {
		categories[v.UsageID] = v.Usage.Challenge.Category
	}
	allCount := make(map[string]int64)
	for _, v := range flags {
		allCount[v.Usage.Challenge.Category] += 1
	}
	solvedCount := make(map[string]int64)
	for _, flag := range solved {
		solvedCount[categories[flag.UsageID]] += 1
	}
	data := make([]gin.H, 0)
	for k, v := range allCount {
		if _, ok := solvedCount[k]; !ok {
			solvedCount[k] = 0
		}
		data = append(data, gin.H{"category": k, "solved": solvedCount[k], "all": v})
	}
	return data
}

func GetTeamRankingResp(teamsData []struct {
	Team   model.Team
	Solved []model.Flag
}, flags []model.Flag, admin bool) gin.H {
	data := make([]gin.H, 0)
	for _, team := range teamsData {
		tmp := gin.H{
			"id":     team.Team.ID,
			"name":   team.Team.Name,
			"desc":   team.Team.Desc,
			"score":  team.Team.Score,
			"avatar": fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(team.Team.Avatar, "/")),
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
