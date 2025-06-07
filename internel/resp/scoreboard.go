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
	Solved []model.ContestFlag
}, flags []model.ContestFlag, admin bool) gin.H {
	data := make([]gin.H, 0)
	for _, team := range teamsData {
		tmp := gin.H{
			"id":     team.Team.ID,
			"name":   team.Team.Name,
			"desc":   team.Team.Desc,
			"score":  team.Team.Score,
			"avatar": fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(string(team.Team.Avatar), "/")),
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
