package middleware

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
)

var (
	ContestIsComing    = ContestStatus(model.ContestIsComing)
	ContestIsRunning   = ContestStatus(model.ContestIsRunning)
	ContestIsNotOver   = ContestStatus(model.ContestIsComing, model.ContestIsRunning)
	ContestIsNotComing = ContestStatus(model.ContestIsRunning, model.ContestIsOver)
)

func ContestStatus(statusL ...string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		contest := GetContest(ctx)
		if slices.Contains(statusL, contest.Status()) {
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": contest.Status(), "data": nil})
		return
	}
}
