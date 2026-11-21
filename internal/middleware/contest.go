package middleware

import (
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/i18n"
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
		resp.AbortJSON(ctx, model.FailRetVal(contest.Status(), "Contest", fmt.Sprintf("contest_id=%d", contest.ID), nil))
		return
	}
}

// CheckLeaderboardVisible 排行榜黑灯期间选手看不到排行榜，后台排行榜接口不受影响。
func CheckLeaderboardVisible(ctx *gin.Context) {
	contest := GetContest(ctx)
	if contest.LeaderboardHidden {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Model.Contest.LeaderboardHidden, "Contest", fmt.Sprintf("contest_id=%d", contest.ID), nil))
		return
	}
	ctx.Next()
}
