package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckChallengeType(t string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if GetContestChallenge(ctx).Type != t {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.InvalidChallengeType, "data": nil})
			return
		}
		ctx.Next()
	}
}

// CheckSolved model.Team 是否完全解决 model.ContestChallenge
func CheckSolved(ctx *gin.Context) {
	team := GetTeam(ctx)
	contestChallenge := GetContestChallenge(ctx)
	if service.CheckIfSolved(db.DB.WithContext(ctx), team, contestChallenge.ContestFlags) {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.AlreadySolved, "data": nil})
		return
	}
	ctx.Next()
}
