package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SubmitFlag(ctx *gin.Context) {
	var form f.SubmitFlagForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
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
		if contestChallenge.Challenge.Type == model.PodsChallengeType && service.CheckIfSolved(db.DB.WithContext(ctx), team, contestChallenge) {
			service.StopVictim(db.DB.WithContext(ctx), team, contestChallenge)
		}
	}(ctx.Copy())
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
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	contestFlag := middleware.GetContestFlag(ctx)
	tx := db.DB.WithContext(ctx).Begin()
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
