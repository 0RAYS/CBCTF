package router

import (
	"os"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
)

func GetChallenge(ctx *gin.Context) {
	challenge, ret := service.GetChallengeWithFlags(db.DB.WithContext(ctx.Request.Context()), middleware.GetChallenge(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(resp.GetChallengeResp(service.BuildChallengeView(db.DB.WithContext(ctx.Request.Context()), challenge))))
}

func GetChallenges(ctx *gin.Context) {
	var form dto.GetChallengesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	challenges, count, ret := service.ListChallengeViews(db.DB.WithContext(ctx.Request.Context()), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, challenge := range challenges {
		data = append(data, resp.GetChallengeResp(challenge))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "challenges": data}))
}

func GetAvailableChallenges(ctx *gin.Context) {
	var form dto.GetChallengesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	challenges, count, ret := service.ListAvailableChallenges(db.DB.WithContext(ctx.Request.Context()), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, challenge := range challenges {
		data = append(data, resp.GetSimpleChallengeResp(challenge))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "challenges": data}))
}

func GetChallengeCategories(ctx *gin.Context) {
	var form dto.GetCategoriesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	categories, ret := service.ListChallengeCategories(db.DB.WithContext(ctx.Request.Context()), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(categories))
}

func CreateChallenge(ctx *gin.Context) {
	var form dto.CreateChallengeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateChallengeEventType)
	challenge, ret := service.CreateChallenge(db.DB.WithContext(ctx.Request.Context()), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if err := os.MkdirAll(challenge.BasicDir(), 0755); err != nil {
		log.Logger.Warningf("create challenge dir err: %s", err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Model.File.CreateDirError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetChallengeResp(service.BuildChallengeView(db.DB.WithContext(ctx.Request.Context()), challenge))))
}

func UpdateChallenge(ctx *gin.Context) {
	var form dto.UpdateChallengeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateChallengeEventType)
	ret := service.UpdateChallenge(db.DB.WithContext(ctx.Request.Context()), middleware.GetChallenge(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, ret)
}

func DeleteChallenge(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteChallengeEventType)
	challenge := middleware.GetChallenge(ctx)
	// 题库题目删除不影响比赛题目副本，这里先统计引用数量，供前端提示。
	usageCount, usageRet := service.CountChallengeUsages(db.DB.WithContext(ctx.Request.Context()), challenge)
	if !usageRet.OK {
		resp.JSON(ctx, usageRet)
		return
	}
	ret := service.DeleteChallenge(db.DB.WithContext(ctx.Request.Context()), challenge)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if err := os.RemoveAll(challenge.BasicDir()); err != nil {
		log.Logger.Warningf("Failed to remove challenge basic dir: %s", err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"contest_challenge_count": usageCount}))
}

// GetChallengeUsages 返回引用了该题库题目的比赛题目实例。
// 用于在修改或删除题库题目前提示其影响范围。
func GetChallengeUsages(ctx *gin.Context) {
	challenge := middleware.GetChallenge(ctx)
	usages, ret := service.ListChallengeUsages(db.DB.WithContext(ctx.Request.Context()), challenge)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(usages))
	for _, usage := range usages {
		data = append(data, resp.GetChallengeUsageResp(usage))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": len(data), "usages": data}))
}
