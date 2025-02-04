package router

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SubmitFlag(ctx *gin.Context) {
	var form constants.SubmitFlagForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	team := middleware.GetTeam(ctx)
	submission, ok, msg := db.CreateSubmission(ctx, contest, team, middleware.GetSelf(ctx).(model.User), middleware.GetChallenge(ctx), form.Flag)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if submission.Solved {
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "FlagNotMatch", "data": nil})
}

func GetSubmissions(ctx *gin.Context) {
	var form constants.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	submissions, count, ok, msg := db.GetSubmissions(ctx, form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"submissions": submissions, "count": count}})
}
