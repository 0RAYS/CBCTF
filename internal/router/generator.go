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

func GetContestGenerators(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	contest := middleware.GetContest(ctx)
	generators, count, ret := db.InitGeneratorRepo(db.DB).List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
	})
	if ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, generator := range generators {
		data = append(data, resp.GetGeneratorResp(generator))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "generators": data}))
}

func StartContestGenerator(ctx *gin.Context) {
	var form dto.StartGeneratorsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.StartGeneratorEventType)
	contest := middleware.GetContest(ctx)
	go service.StartContestGenerators(db.DB, contest, form)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

func StopContestGenerator(ctx *gin.Context) {
	var form dto.StopGeneratorsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.StopGeneratorEventType)
	go service.StopContestGenerators(db.DB, form)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
