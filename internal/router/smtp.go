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

func GetSmtps(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	smtps, count, ret := service.ListSmtps(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, smtp := range smtps {
		data = append(data, resp.GetSmtpResp(smtp))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"smtps": data, "count": count}))
}

func CreateSmtp(ctx *gin.Context) {
	var form dto.CreateSmtpForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateSmtpEventType)
	smtp, ret := service.CreateSmtp(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetSmtpResp(smtp)))
}

func UpdateSmtp(ctx *gin.Context) {
	var form dto.UpdateSmtpForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateSmtpEventType)
	if _, ret := service.UpdateSmtp(db.DB, middleware.GetSmtp(ctx), form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

func DeleteSmtp(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteSmtpEventType)
	if ret := service.DeleteSmtp(db.DB, middleware.GetSmtp(ctx)); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
