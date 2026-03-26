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
	challenges, count, ret := service.GetChallenges(db.DB, form)
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
	contest := middleware.GetContest(ctx)
	challenges, count, ret := db.InitChallengeRepo(db.DB).ListChallengesNotInContest(contest.ID,
		form.Limit, form.Offset, form.Name, form.Description, form.Category, form.Type)
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
	categories, ret := db.InitChallengeRepo(db.DB).ListCategories(form.Type)
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
	tx := db.DB.Begin()
	challenge, ret := service.CreateChallenge(tx, form)
	if !ret.OK {
		tx.Rollback()
		resp.JSON(ctx, ret)
		return
	}
	tx.Commit()
	if err := os.MkdirAll(challenge.BasicDir(), 0755); err != nil {
		log.Logger.Warningf("create challenge dir err: %s", err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Model.File.CreateDirError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetChallengeResp(challenge)))
}

func UpdateChallenge(ctx *gin.Context) {
	var form dto.UpdateChallengeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateChallengeEventType)
	var ret model.RetVal
	challenge := middleware.GetChallenge(ctx)
	challenge.ChallengeFlags, _, ret = db.InitChallengeFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"challenge_id": challenge.ID},
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	tx := db.DB.Begin()
	ret = service.UpdateChallenge(tx, challenge, form)
	if !ret.OK {
		tx.Rollback()
		resp.JSON(ctx, ret)
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, ret)
}

func DeleteChallenge(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteChallengeEventType)
	challenge := middleware.GetChallenge(ctx)
	tx := db.DB.Begin()
	ret := db.InitChallengeRepo(tx).Delete(challenge.RandID)
	if !ret.OK {
		tx.Rollback()
		resp.JSON(ctx, ret)
		return
	}
	tx.Commit()
	if err := os.RemoveAll(challenge.BasicDir()); err != nil {
		log.Logger.Warningf("Failed to remove challenge basic dir: %s", err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
