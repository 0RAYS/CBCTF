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
	contestFlags, _, ok, msg := db.InitContestFlagRepo(db.DB.WithContext(ctx)).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if !service.CheckIfGenerated(db.DB.WithContext(ctx), team, contestFlags) {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.TeamFlagNotFound, "data": nil})
		return
	}
	ctx.Next()
}
