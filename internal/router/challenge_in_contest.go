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

// GetChallengesInContest 返回比赛对选手可见的题目。
func GetChallengesInContest(ctx *gin.Context) {
	var form dto.GetChallengesInContestForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	challenges, count, ret := service.ListChallengesInContest(db.DB.WithContext(ctx.Request.Context()), middleware.GetContest(ctx), middleware.GetTeam(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(challenges))
	for _, challenge := range challenges {
		data = append(data, resp.GetChallengeRuntimeResp(challenge))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"challenges": data, "count": count}))
}

// GetChallengeCategoriesInContest 返回比赛内的题目分类。
func GetChallengeCategoriesInContest(ctx *gin.Context) {
	var form dto.GetCategoriesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	categories, ret := service.ListChallengeCategoriesInContest(db.DB.WithContext(ctx.Request.Context()), middleware.GetContest(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(categories))
}

// GetChallengeStatus 返回队伍在当前题目上的状态。
func GetChallengeStatus(ctx *gin.Context) {
	status, ret := service.GetChallengeStatus(
		db.DB.WithContext(ctx.Request.Context()),
		middleware.GetTeam(ctx),
		middleware.GetChallenge(ctx),
	)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(resp.GetChallengeStatusResp(status)))
}

// GetAllChallengesInContest 返回比赛下的全部题目，供管理员管理。
func GetAllChallengesInContest(ctx *gin.Context) {
	var form dto.GetAllChallengesInContestForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	challenges, count, ret := service.ListAllChallengesInContest(db.DB.WithContext(ctx.Request.Context()), middleware.GetContest(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(challenges))
	for _, challenge := range challenges {
		data = append(data, resp.GetAdminChallengeResp(challenge))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"challenges": data, "count": count}))
}

// CopyChallengesToContest 把已有题目复制到本比赛，副本与源题互不关联。
func CopyChallengesToContest(ctx *gin.Context) {
	var form dto.CopyChallengesToContestForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CopyChallengeEventType)
	challenges, failed, ret := service.CopyChallengesToContest(db.DB.WithContext(ctx.Request.Context()), middleware.GetContest(ctx), form)
	data := make([]gin.H, 0, len(challenges))
	for _, challenge := range challenges {
		data = append(data, resp.GetAdminChallengeResp(challenge))
	}
	ctx.Set(middleware.CTXEventSuccessKey, ret.OK)
	ret.Data = gin.H{"challenges": data, "failed": failed, "batch": ret.Data}
	resp.JSON(ctx, ret)
}
