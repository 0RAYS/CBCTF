package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetTeamResp(team model.Team, admin bool) gin.H {
	data := gin.H{
		"id":         team.ID,
		"name":       team.Name,
		"score":      team.Score,
		"avatar":     team.Avatar,
		"last":       team.Last,
		"users":      len(team.Users),
		"desc":       team.Desc,
		"captain_id": team.CaptainID,
	}
	if admin {
		data["captcha"] = team.Captcha
	}
	return data
}
