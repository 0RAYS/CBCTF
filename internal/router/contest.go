package router

import (
	"CBCTF/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetContest(ctx *gin.Context) {
	type contestIDUri struct {
		ContestID uint `uri:"contestID" binding:"required"`
	}
	var uri contestIDUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest, ok, msg := db.GetContestByID(ctx, uri.ContestID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": contest})
}

func GetContests(ctx *gin.Context) {
	var form GetModelsForm
	self, _ := ctx.Get("Self")
	all := false
	if self.(map[string]interface{})["Type"].(string) == "admin" {
		all = true
	}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contests, count, ok, msg := db.GetContests(ctx, form.Limit, form.Offset, all)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"count": count, "contests": contests}})
}
