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

func GetNotices(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	notices, count, ret := service.ListNotices(db.DB, middleware.GetContest(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, notice := range notices {
		data = append(data, resp.GetNoticeResp(notice))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "notices": data}))
}

func CreateNotice(ctx *gin.Context) {
	var form dto.CreateNoticeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateNoticeEventType)
	notice, ret := service.CreateNotice(db.DB, middleware.GetContest(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(notice))
}

func UpdateNotice(ctx *gin.Context) {
	var form dto.UpdateNoticeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateNoticeEventType)
	notice := middleware.GetNotice(ctx)
	ret := service.UpdateNotice(db.DB, notice, form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func DeleteNotice(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteNoticeEventType)
	notice := middleware.GetNotice(ctx)
	ret := service.DeleteNotice(db.DB, notice)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}
