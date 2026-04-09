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

func GetGenerators(ctx *gin.Context) {
	var form dto.ListGeneratorsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	generators, count, ret := service.ListGenerators(db.DB, middleware.GetContest(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, generator := range generators {
		data = append(data, resp.GetGeneratorResp(generator))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "generators": data}))
}

func StartGenerator(ctx *gin.Context) {
	var form dto.StartGeneratorsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.StartGeneratorEventType)
	contest := middleware.GetContest(ctx)
	go service.StartGenerators(db.DB, contest.ID, form)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

func StopGenerator(ctx *gin.Context) {
	var form dto.StopGeneratorsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.StopGeneratorEventType)
	go service.StopGenerators(db.DB, form)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
