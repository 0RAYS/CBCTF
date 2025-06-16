package middleware

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ContestIsNotOver   = ContestStatus(model.ContestIsComing, model.ContestIsRunning)
	ContestIsNotComing = ContestStatus(model.ContestIsRunning, model.ContestIsOver)
)

func ContestStatus(statusL ...string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		contest := GetContest(ctx)
		for _, status := range statusL {
			if contest.Status() == status {
				ctx.Next()
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": contest.Status(), "data": nil})
		return
	}
}
