package middleware

import (
	"CBCTF/internel/i18n"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CheckSolved model.Team 是否完全解决 model.Usage
func CheckSolved(ctx *gin.Context) {
	team := GetTeam(ctx)
	contestChallenge := GetContestChallenge(ctx)
	if service.CheckIfSolved(db.DB.WithContext(ctx), team, contestChallenge) {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.AlreadySolved, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Next()
}
