package router

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
)

// GetRangePools 列出比赛里靶场题目的实例池状态。
func GetRangePools(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	pools, ret := service.BuildRangePoolViews(db.DB.WithContext(ctx.Request.Context()), contest.ID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(resp.GetRangePoolsResp(pools)))
}

// PrewarmRangePools 手动补齐比赛里靶场题目的实例池。
// 与定时任务走同一条补货路径,池的常驻数量仍由出题人声明决定。
func PrewarmRangePools(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	if ret := service.EnsureContestRangePools(ctx.Request.Context(), contest.ID); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	pools, ret := service.BuildRangePoolViews(db.DB.WithContext(ctx.Request.Context()), contest.ID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(resp.GetRangePoolsResp(pools)))
}
