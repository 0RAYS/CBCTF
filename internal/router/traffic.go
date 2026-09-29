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
	// Init enables ContextWithFallback so analysis follows request cancellation.
	data, ret := service.GetTraffic(ctx, victim, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func GetTrafficAnalysis(ctx *gin.Context) {
	data, ret := service.GetTrafficAnalysis(ctx, middleware.GetVictim(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func GetContestTrafficOverlaps(ctx *gin.Context) {
	data, ret := service.GetContestTrafficOverlaps(ctx, middleware.GetContest(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}
