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

func GetChallenge(ctx *gin.Context) {
	challenge, ret := service.GetChallengeWithFlags(db.DB.WithContext(ctx.Request.Context()), middleware.GetChallenge(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(resp.GetChallengeResp(service.BuildChallengeView(db.DB.WithContext(ctx.Request.Context()), challenge))))
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
		resp.JSON(ctx, model.FailRetVal(i18n.Model.File.CreateDirError, "Challenge", challenge.BasicDir(), err))
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
	ret := service.DeleteChallenge(db.DB.WithContext(ctx.Request.Context()), challenge)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if err := os.RemoveAll(challenge.BasicDir()); err != nil {
		log.Logger.Warningf("Failed to remove challenge basic dir: %s", err)
		resp.JSON(ctx, model.FailRetVal(i18n.Common.UnknownError, "Challenge", challenge.BasicDir(), err))
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

// CopyChallenge 把题目复制成目标比赛下的一道新题，副本与源题互不关联。
func CopyChallenge(ctx *gin.Context) {
	var form dto.CopyChallengeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CopyChallengeEventType)
	source, ret := service.GetChallengeWithFlags(db.DB.WithContext(ctx.Request.Context()), middleware.GetChallenge(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	challenge, ret := service.CopyChallenge(db.DB.WithContext(ctx.Request.Context()), source, form.ContestID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetChallengeResp(service.BuildChallengeView(db.DB.WithContext(ctx.Request.Context()), challenge))))
}

// ResetChallengeSolves 清空题目的解题记录并重算分数与排名。
func ResetChallengeSolves(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.ResetChallengeSolvesEventType)
	ret := service.ResetChallengeSolves(db.DB.WithContext(ctx.Request.Context()), middleware.GetChallenge(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, ret)
}
