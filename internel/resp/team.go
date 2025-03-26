package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetTeamResp(team model.Team, solved []model.Flag, flags []model.Flag) gin.H {
	data := GetSolvedStateResp(solved, flags)
	data["name"] = team.Name
	data["score"] = team.Score
	data["avatar"] = team.Avatar
	data["last"] = team.Last
	data["users"] = len(team.Users)
	data["desc"] = team.Desc
	data["contest_id"] = team.ContestID
	data["captain_id"] = team.CaptainID
	return data
}
