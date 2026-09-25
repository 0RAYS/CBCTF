package middleware

import (
	"slices"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
	"CBCTF/internal/resp"
)

var (
	ContestIsComing    = ContestStatus(model.ContestIsComing)
	ContestIsRunning   = ContestStatus(model.ContestIsRunning)
	ContestIsNotOver   = ContestStatus(model.ContestIsComing, model.ContestIsRunning)
	ContestIsNotComing = ContestStatus(model.ContestIsRunning, model.ContestIsOver)
)

func ContestStatus(statusL ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		contest := GetContest(ctx)
		if slices.Contains(statusL, contest.Status()) {
			ctx.Next()
			return
		}
		resp.AbortJSON(ctx, model.RetVal{Msg: contest.Status()})
		return
	}
}
