package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetTeamResp(team model.Team) gin.H {
	return gin.H{
		"name":       team.Name,
		"score":      team.Score,
		"avatar":     team.Avatar,
		"last":       team.Last,
		"users":      len(team.Users),
		"desc":       team.Desc,
		"captain_id": team.CaptainID,
	}
}
