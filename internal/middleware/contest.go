package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckRunning(ctx *gin.Context) {
	contest := GetContest(ctx)
	if !contest.IsRunning() {
		ctx.JSON(http.StatusOK, gin.H{"msg": contest.Status(), "data": nil})
		ctx.Abort()
		return
	}
	ctx.Next()
}
