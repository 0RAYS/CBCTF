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

func SubmitFlag(ctx *gin.Context) {
	var form dto.SubmitFlagForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.SubmitFlagEventType)
	ret := service.SubmitContestFlag(
		db.DB,
		middleware.GetSelf(ctx),
		middleware.GetTeam(ctx),
		middleware.GetContest(ctx),
		middleware.GetChallenge(ctx),
		middleware.GetContestChallenge(ctx),
		form,
		ctx.ClientIP(),
	)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, ret)
}

func GetContestFlags(ctx *gin.Context) {
	contestFlags, ret := service.ListContestFlags(db.DB, middleware.GetContestChallenge(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(contestFlags))
	for _, contestFlag := range contestFlags {
		data = append(data, resp.GetContestFlagResp(contestFlag))
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func UpdateContestFlag(ctx *gin.Context) {
	var form dto.UpdateContestFlagForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateContestChallengeFlagEventType)
	ret := service.UpdateContestFlag(db.DB, middleware.GetContestChallenge(ctx), middleware.GetContestFlag(ctx), form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}
