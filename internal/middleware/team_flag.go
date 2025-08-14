package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
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
