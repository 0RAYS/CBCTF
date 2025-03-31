package middleware

import (
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"CBCTF/internel/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CheckGenerated model.Team 是否初始化 model.Usage
func CheckGenerated(ctx *gin.Context) {
	usage := GetUsage(ctx)
	team := GetTeam(ctx)
	if !service.IsGenerated(db.DB.WithContext(ctx), usage, team) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "AnswerNotFound", "data": nil})
		ctx.Abort()
		return
	}
	ctx.Next()
}

// CheckSolved model.Team 是否完全解决 model.Usage
func CheckSolved(ctx *gin.Context) {
	usage := GetUsage(ctx)
	team := GetTeam(ctx)
	flags, ok, msg := service.GetTeamSolved(db.DB.WithContext(ctx), team.ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	solved := make([]uint, 0)
	for _, f := range flags {
		solved = append(solved, f.ID)
	}
	for _, f := range usage.Flags {
		if !utils.In(f.ID, solved) {
			ctx.Next()
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "AlreadySolved", "data": nil})
	ctx.Abort()
}
