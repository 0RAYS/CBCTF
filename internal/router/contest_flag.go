package router

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	db "CBCTF/internal/repo"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SubmitFlag(ctx *gin.Context) {
	var form f.SubmitFlagForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.SubmitFlagEventType)
	user := middleware.GetSelf(ctx).(model.User)
	team := middleware.GetTeam(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	result, _, ok, msg := service.Submit(tx, user, team, contestChallenge, form, ctx.ClientIP())
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	go func(ctx *gin.Context) {
		if contestChallenge.Type == model.PodsChallengeType && service.CheckIfSolved(db.DB.WithContext(ctx), team, contestChallenge) {
			service.StopTeamVictim(db.DB.WithContext(ctx), team, contestChallenge)
		}
	}(ctx.Copy())
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": result, "data": nil})
}

func GetContestFlags(ctx *gin.Context) {
	contestChallenge := middleware.GetContestChallenge(ctx)
	data := make([]gin.H, 0)
	for _, contestFlag := range contestChallenge.ContestFlags {
		data = append(data, resp.GetContestFlagResp(contestFlag))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
}

func GetContestFlag(ctx *gin.Context) {
	contestFlag := middleware.GetContestFlag(ctx)
	data := resp.GetContestFlagResp(contestFlag)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
}

func UpdateContestFlag(ctx *gin.Context) {
	var form f.UpdateContestFlagForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	contestChallenge := middleware.GetContestChallenge(ctx)
	contestFlag := middleware.GetContestFlag(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	if contestChallenge.Type == model.QuestionChallengeType && form.Value != nil {
		form.Value = &contestFlag.Value
	}
	ok, msg := db.InitContestFlagRepo(tx).Update(contestFlag.ID, db.UpdateContestFlagOptions{
		Value:     form.Value,
		Score:     form.Score,
		Decay:     form.Decay,
		MinScore:  form.MinScore,
		ScoreType: form.ScoreType,
	})
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
