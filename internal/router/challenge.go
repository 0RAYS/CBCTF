package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"os"

	"github.com/gin-gonic/gin"
)

func GetChallenges(ctx *gin.Context) {
	var form dto.GetChallengesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	challenges, count, ret := service.ListChallengeViews(db.DB, form)
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

func GetChallengeNotInContest(ctx *gin.Context) {
	var form dto.GetChallengesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	challenges, count, ret := service.ListChallengesNotInContest(db.DB, middleware.GetContest(ctx), form)
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
	categories, ret := service.ListChallengeCategories(db.DB, form)
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
	challenge, ret := service.CreateChallengeWithTransaction(db.DB, form)
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
	resp.JSON(ctx, model.SuccessRetVal(resp.GetChallengeResp(service.GetChallengeView(db.DB, challenge))))
}

func UpdateChallenge(ctx *gin.Context) {
	var form dto.UpdateChallengeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateChallengeEventType)
	ret := service.UpdateChallengeWithTransaction(db.DB, middleware.GetChallenge(ctx), form)
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
	ret := service.DeleteChallengeWithTransaction(db.DB, challenge)
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
	resp.JSON(ctx, model.SuccessRetVal())
}
