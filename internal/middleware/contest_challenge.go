package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CheckSolved model.Team 是否完全解决 model.Usage
func CheckSolved(ctx *gin.Context) {
	team := GetTeam(ctx)
	contestChallenge := GetContestChallenge(ctx)
	if service.CheckIfSolved(db.DB.WithContext(ctx), team, contestChallenge) {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.AlreadySolved, "data": nil})
		return
	}
	ctx.Next()
}
