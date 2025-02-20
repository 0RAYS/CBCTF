package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SubmitFlag(ctx *gin.Context) {
	var form f.SubmitFlagForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	team := middleware.GetTeam(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	submission, ok, msg := db.CreateSubmission(tx, contest, team, middleware.GetSelf(ctx).(model.User), middleware.GetChallenge(ctx), form.Flag)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	if submission.Solved {
		go db.UpdateRanking(db.DB, contest.ID)
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "FlagNotMatch", "data": nil})
}

func GetSubmissions(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	submissions, count, ok, msg := db.GetSubmissions(db.DB.WithContext(ctx), form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"submissions": &submissions, "count": count}})
}

func GetTeamSubmissions(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	submissions, count, ok, msg := db.GetSubmissions(db.DB.WithContext(ctx), form.Limit, form.Offset, middleware.GetTeam(ctx).ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"submissions": &submissions, "count": count}})
}
