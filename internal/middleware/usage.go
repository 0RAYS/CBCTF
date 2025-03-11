package middleware

import (
	"CBCTF/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
