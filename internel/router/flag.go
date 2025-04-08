package router

import (
	f "CBCTF/internel/form"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	user := middleware.GetSelf(ctx).(model.User)
	team := middleware.GetTeam(ctx)
	usage := middleware.GetUsage(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	submission, ok, msg := service.Submit(tx, user, team, usage, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	if submission.Solved {
		//TODO 使用定时任务代替, 提升性能
		//go service.UpdateTeamRanking(db.DB.WithContext(ctx), team.ContestID)
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "FlagNotMatch", "data": nil})
}

func GetFlag(ctx *gin.Context) {
	flag := middleware.GetFlag(ctx)
	data := resp.GetFlagResp(flag)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}

func GetFlags(ctx *gin.Context) {
	usage := middleware.GetUsage(ctx)
	repo := db.InitFlagRepo(db.DB.WithContext(ctx))
	flags, _, ok, msg := repo.GetByKeyID("usage_id", usage.ID, -1, -1)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, flag := range flags {
		data = append(data, resp.GetFlagResp(flag))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}

func UpdateFlag(ctx *gin.Context) {
	var form f.UpdateFlagForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	flag := middleware.GetFlag(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.UpdateFlag(tx, flag, form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
