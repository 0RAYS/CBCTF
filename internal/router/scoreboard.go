package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRank(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBindQuery(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	teams, ok, msg := db.GetRanking(contest.ID, form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": teams})
}
