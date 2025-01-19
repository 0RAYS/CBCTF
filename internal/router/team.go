package router

import (
	"github.com/gin-gonic/gin"
)

func UpdateTeam(ctx *gin.Context) {
	//self, _ := ctx.Get("Self")
	//var form UpdateTeamForm
	//type teamIDUri struct {
	//	TeamID uint `uri:"teamID" binding:"required"`
	//}
	//var uri teamIDUri
	//if err := ctx.ShouldBindUri(&uri); err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
	//	return
	//}
	//if err := ctx.ShouldBindJSON(&form); err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
	//	return
	//}
	//team, ok, msg := db.GetTeamByID(ctx, uri.TeamID, false)
	//if !ok {
	//	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
	//	return
	//}
	//data := utils.Form2Map(form)
	//if team.Name != data["name"].(string) && !db.IsUniqueTeamName(data["name"].(string), team.ContestID) {
	//	ctx.JSON(http.StatusOK, gin.H{"msg": "TeamNameExists", "data": nil})
	//	return
	//}
	//if team.CaptainID != data["captain_id"].(uint) {
	//
	//}
	//ok, msg := db.UpdateTeam(ctx, uri.TeamID)
	//if !ok {
	//	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
	//} else {
	//	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
	//}
}
