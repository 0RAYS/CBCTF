package middleware

import (
	"CBCTF/internel/i18n"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CheckIfGenerated model.Team 是否初始化 model.TeamFlag
func CheckIfGenerated(ctx *gin.Context) {
	contestChallenge := GetContestChallenge(ctx)
	team := GetTeam(ctx)
	if !service.CheckIfGenerated(db.DB.WithContext(ctx), team, contestChallenge) {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.TeamFlagNotFound, "data": nil})
		return
	}
	ctx.Next()
}
