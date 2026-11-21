package router

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
)

func GetContest(ctx *gin.Context) {
	contestView := service.BuildContestView(db.DB.WithContext(ctx.Request.Context()), middleware.GetContest(ctx))
	resp.JSON(ctx, model.SuccessRetVal(resp.GetContestResp(contestView, middleware.IsFullAccess(ctx))))
}

func GetContests(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	contests, count, ret := service.ListContests(db.DB.WithContext(ctx.Request.Context()), form, middleware.IsFullAccess(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, contest := range contests {
		data = append(data, resp.GetContestResp(contest, middleware.IsFullAccess(ctx)))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"contests": data, "count": count}))
}

func CreateContest(ctx *gin.Context) {
	var form dto.CreateContestForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateContestEventType)
	contest, ret := service.CreateContest(db.DB.WithContext(ctx.Request.Context()), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetContestResp(service.BuildContestView(db.DB.WithContext(ctx.Request.Context()), contest), true)))
}

func UpdateContest(ctx *gin.Context) {
	var form dto.UpdateContestForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateContestEventType)
	contest := middleware.GetContest(ctx)
	ret := service.UpdateContest(db.DB.WithContext(ctx.Request.Context()), contest, form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func DeleteContest(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteContestEventType)
	ret := service.DeleteContest(db.DB.WithContext(ctx.Request.Context()), middleware.GetContest(ctx))
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

// UpdateContestLeaderboard 开启或关闭排行榜黑灯模式。
func UpdateContestLeaderboard(ctx *gin.Context) {
	var form dto.UpdateContestLeaderboardForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateContestLeaderboardEventType)
	contest, ret := service.UpdateContestLeaderboard(db.DB.WithContext(ctx.Request.Context()), middleware.GetContest(ctx), form.Hidden)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	view := service.BuildContestView(db.DB.WithContext(ctx.Request.Context()), contest)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetContestResp(view, true)))
}

// UpdateContestPause 暂停或恢复比赛。
func UpdateContestPause(ctx *gin.Context) {
	var form dto.UpdateContestPauseForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateContestPauseEventType)
	contest, ret := service.UpdateContestPause(db.DB.WithContext(ctx.Request.Context()), middleware.GetContest(ctx), form.Paused)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	view := service.BuildContestView(db.DB.WithContext(ctx.Request.Context()), contest)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetContestResp(view, true)))
}
