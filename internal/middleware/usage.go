package middleware

import (
	"CBCTF/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CheckSolved model.Team 是否解决过 model.Usage
func CheckSolved(ctx *gin.Context) {
	usage := GetUsage(ctx)
	team := GetTeam(ctx)
	if db.IsSolved(db.DB.WithContext(ctx), usage.ContestID, team.ID, usage.ChallengeID) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "AlreadySolved", "data": nil})
		ctx.Abort()
		return
	}
	ctx.Next()
}
