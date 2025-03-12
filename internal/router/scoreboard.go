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
	teams, count, ok, msg := db.GetTeamRanking(db.DB.WithContext(ctx), contest.ID, form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	var tmp []gin.H
	for _, team := range teams {
		state, _, _ := db.GetTeamSolvedState(db.DB.WithContext(ctx), team)
		tmp = append(tmp, gin.H{
			"team":   &team,
			"solved": state,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"count": count, "teams": tmp}})
}

func GetRankDetail(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBindQuery(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	data, ok, msg := db.GetTeamRankDetail(db.DB.WithContext(ctx), contest.ID, form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}
