package router

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SubmitFlag(ctx *gin.Context) {
	var form constants.SubmitFlagForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest, ok, msg := db.GetContestByID(ctx, middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	team, ok, msg := db.GetTeamByUserID(ctx, middleware.GetSelfID(ctx), middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	submission, ok, msg := db.CreateSubmission(ctx, contest.ID, team.ID, middleware.GetSelfID(ctx), middleware.GetChallengeID(ctx), form.Flag)
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
