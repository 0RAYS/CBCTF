package router

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
)

func GetTraffics(ctx *gin.Context) {
	var form dto.GetTrafficForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	victim := middleware.GetVictim(ctx)
	data, ret := service.GetTraffic(ctx.Request.Context(), victim, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func GetTrafficAnalysis(ctx *gin.Context) {
	data, ret := service.GetTrafficAnalysis(ctx.Request.Context(), middleware.GetVictim(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func GetContestTrafficOverlaps(ctx *gin.Context) {
	data, ret := service.GetContestTrafficOverlaps(ctx.Request.Context(), middleware.GetContest(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}
