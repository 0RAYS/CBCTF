package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckTeamVictimCount(ctx *gin.Context) {
	contest := GetContest(ctx)
	team := GetTeam(ctx)
	count, ok, msg := service.CountTeamVictims(db.DB, team)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if count >= contest.Victims {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.VictimLimited, "data": nil})
		return
	}
	ctx.Next()
}
