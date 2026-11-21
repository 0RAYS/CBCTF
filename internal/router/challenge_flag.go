package router

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
)

func GetChallengeFlags(ctx *gin.Context) {
	flags, ret := service.ListChallengeFlags(db.DB.WithContext(ctx.Request.Context()), middleware.GetChallenge(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(flags))
	for _, flag := range flags {
		data = append(data, resp.GetChallengeFlagResp(flag))
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func GetChallengeFlag(ctx *gin.Context) {
	resp.JSON(ctx, model.SuccessRetVal(resp.GetChallengeFlagResp(middleware.GetChallengeFlag(ctx))))
}

func UpdateChallengeFlag(ctx *gin.Context) {
	var form dto.UpdateChallengeFlagForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateChallengeFlagEventType)
	ret := service.UpdateChallengeFlag(db.DB.WithContext(ctx.Request.Context()), middleware.GetChallengeFlag(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, ret)
}

func GetChallengeFlagSolvers(ctx *gin.Context) {
	solvers, ret := service.ListChallengeFlagSolvers(db.DB.WithContext(ctx.Request.Context()), middleware.GetChallengeFlag(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(solvers))
	for _, solver := range solvers {
		data = append(data, resp.GetFlagSolverResp(solver))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"solvers": data, "count": int64(len(data))}))
}

func SubmitFlag(ctx *gin.Context) {
	var form dto.SubmitFlagForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.SubmitFlagEventType)
	ret := service.SubmitChallengeFlag(
		db.DB.WithContext(ctx.Request.Context()),
		middleware.GetSelf(ctx),
		middleware.GetTeam(ctx),
		middleware.GetContest(ctx),
		middleware.GetChallenge(ctx),
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
