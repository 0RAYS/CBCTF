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
	options := db.GetOptions{
		Deleted: form.Deleted,
		Sort:    []string{"id DESC"},
	}
	contest := middleware.GetContest(ctx)
	if contest.ID > 0 {
		options.Conditions = map[string]any{"contest_id": contest.ID}
	}
	generators, count, ret := db.InitGeneratorRepo(db.DB).List(form.Limit, form.Offset, options)
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
