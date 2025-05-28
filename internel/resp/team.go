package resp

import (
	"CBCTF/internel/config"
	"CBCTF/internel/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetTeamResp(team model.Team) gin.H {
	data := gin.H{
		"id":         team.ID,
		"name":       team.Name,
		"score":      team.Score,
		"avatar":     fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(team.Avatar, "/")),
		"last":       team.Last,
		"rank":       team.Rank,
		"users":      len(team.Users),
		"desc":       team.Desc,
		"captain_id": team.CaptainID,
	}
	return data
}
