package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"

	"github.com/gin-gonic/gin"
)

func GetContest(ctx *gin.Context) {
	contestView := service.GetContestView(db.DB, middleware.GetContest(ctx))
	resp.JSON(ctx, model.SuccessRetVal(resp.GetContestResp(contestView, middleware.IsFullAccess(ctx))))
}

func GetContests(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	contests, count, ret := service.ListContests(db.DB, form, middleware.IsFullAccess(ctx))
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
	contest, ret := service.CreateContest(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetContestResp(service.GetContestView(db.DB, contest), true)))
}

func UpdateContest(ctx *gin.Context) {
	var form dto.UpdateContestForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateContestEventType)
	contest := middleware.GetContest(ctx)
	ret := service.UpdateContest(db.DB, contest, form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func DeleteContest(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteContestEventType)
	ret := service.DeleteContestWithTransaction(db.DB, middleware.GetContest(ctx))
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}
