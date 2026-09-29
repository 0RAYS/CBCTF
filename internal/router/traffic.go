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
	// Gin's context does not forward cancellation unless ContextWithFallback is enabled.
	data, ret := service.GetTraffic(ctx.Request.Context(), victim, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}
