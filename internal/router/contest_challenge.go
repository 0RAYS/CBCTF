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

func GetContestChallenges(ctx *gin.Context) {
	var form dto.GetContestChallengesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	challenges, count, ret := service.ListContestChallengeViews(db.DB, middleware.GetContest(ctx), middleware.GetTeam(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(challenges))
	for _, challenge := range challenges {
		item := resp.GetContestChallengeResp(challenge)
		item["hidden"] = false
		data = append(data, item)
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"challenges": data, "count": count}))
}

func GetAllContestChallenges(ctx *gin.Context) {
	var form dto.GetAllContestChallengesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	contestChallenges, count, ret := service.ListAdminContestChallenges(db.DB, middleware.GetContest(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(contestChallenges))
	for _, contestChallenge := range contestChallenges {
		data = append(data, resp.GetAdminContestChallengeResp(contestChallenge))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"challenges": data, "count": count}))
}

func GetContestChallengeCategories(ctx *gin.Context) {
	var form dto.GetCategoriesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	categories, ret := service.ListContestChallengeCategories(db.DB, middleware.GetContest(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(categories))
}

func GetContestChallengeStatus(ctx *gin.Context) {
	status, ret := service.GetContestChallengeStatus(
		db.DB,
		middleware.GetTeam(ctx),
		middleware.GetChallenge(ctx),
		middleware.GetContestChallenge(ctx),
	)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(resp.GetContestChallengeStatusResp(status)))
}

func AddContestChallenge(ctx *gin.Context) {
	var form dto.CreateContestChallengeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateContestChallengeEventType)
	contestChallenges, failed, _ := service.CreateContestChallenge(db.DB, middleware.GetContest(ctx), form)
	data := make([]gin.H, 0, len(contestChallenges))
	for _, contestChallenge := range contestChallenges {
		data = append(data, resp.GetAdminContestChallengeResp(contestChallenge))
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"contest_challenge": data, "failed": failed}))
}

func UpdateContestChallenge(ctx *gin.Context) {
	var form dto.UpdateContestChallengeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateContestChallengeEventType)
	ret := service.UpdateContestChallenge(db.DB, middleware.GetContestChallenge(ctx), form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func DeleteContestChallenge(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteContestChallengeEventType)
	ret := service.DeleteContestChallenge(db.DB, middleware.GetContestChallenge(ctx))
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func GetContestFlagSolvers(ctx *gin.Context) {
	solvers, ret := service.ListContestFlagSolvers(db.DB, middleware.GetContestFlag(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(solvers))
	for _, solver := range solvers {
		data = append(data, resp.GetContestFlagSolverResp(solver))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"solvers": data, "count": int64(len(data))}))
}
