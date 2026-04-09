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

func GetBranding(ctx *gin.Context) {
	branding, ret := service.GetDefaultBranding(db.DB)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(resp.GetBrandingResp(branding)))
}

func GetAdminBranding(ctx *gin.Context) {
	GetBranding(ctx)
}

func UpdateBranding(ctx *gin.Context) {
	var form dto.UpdateBrandingForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateBrandingEventType)
	branding, ret := service.UpdateBranding(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetBrandingResp(branding)))
}
