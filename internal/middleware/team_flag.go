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
	team := GetTeam(ctx)
	contestChallenge := GetContestChallenge(ctx)
	if !service.CheckIfGenerated(db.DB.WithContext(ctx), team, contestChallenge.ContestFlags) {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.TeamFlagNotFound, "data": nil})
		return
	}
	ctx.Next()
}
