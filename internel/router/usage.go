package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/k8s"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	"CBCTF/internel/redis"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func GetUsageStatus(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	usage := middleware.GetUsage(ctx)
	DB := db.DB.WithContext(ctx)
	var data gin.H
	data["attempts"] = service.CountAttempts(DB, team, usage)
	data["init"] = service.IsGenerated(DB, usage, team)
	data["solved"] = service.IsSolved(DB, team, usage)
	data["remote"] = service.GetRemoteStatus(DB, usage)
	data["file"] = func() string {
		if _, err := os.Stat(usage.Challenge.AttachmentPath(team.ID)); err != nil {
			return ""
		}
		return usage.Challenge.AttachmentPath(team.ID)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}

func GetUsage(ctx *gin.Context) {
	usage := middleware.GetUsage(ctx)
	data := resp.GetUsageResp(usage, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}

func GetUsages(ctx *gin.Context) {
	var (
		all     = middleware.GetRole(ctx) == "admin"
		DB      = db.DB.WithContext(ctx)
		contest = middleware.GetContest(ctx)
	)
	usages, _, ok, msg := db.InitUsageRepo(DB).GetAll(contest.ID, -1, -1, true, 3, all)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, usage := range usages {
		tmp := resp.GetUsageResp(usage, all)
		if !all {
			team := middleware.GetTeam(ctx)
			tmp["attempts"] = service.CountAttempts(DB, team, usage)
			tmp["init"] = service.IsGenerated(DB, usage, team)
			tmp["solved"] = service.IsSolved(DB, team, usage)
			tmp["remote"] = service.GetRemoteStatus(DB, usage)
			tmp["file"] = func() string {
				if _, err := os.Stat(usage.Challenge.AttachmentPath(team.ID)); err != nil {
					return ""
				}
				return usage.Challenge.AttachmentPath(team.ID)
			}
		}
		data = append(data, tmp)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}

func InitUsage(reset bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var (
			team  = middleware.GetTeam(ctx)
			usage = middleware.GetUsage(ctx)
			tx    = db.DB.WithContext(ctx).Begin()
		)
		if ok, err := redis.CheckChallengeInit(team.ID, usage.ChallengeID); ok || err != nil {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"msg": "TooQuick", "data": nil})
			return
		}
		_ = redis.RecordChallengeInit(team.ID, usage.ChallengeID)
		answers, ok, msg := service.GenerateAnswer(tx, usage, team, reset)
		if !ok {
			tx.Rollback()
		}
		if usage.Challenge.Type == model.DynamicChallenge {
			ok, msg = k8s.GenerateAttachment(usage, team, answers)
		}
		if !ok {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
	}
}

func AddUsage(ctx *gin.Context) {
	var form f.CreateUsageForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	DB := db.DB.WithContext(ctx)
	usages, failed, _, _ := service.CreateUsage(DB, middleware.GetContest(ctx), form)
	data := make([]gin.H, 0)
	for _, usage := range usages {
		data = append(data, resp.GetUsageResp(usage, true))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"usages": data, "failed": failed}})
}

func UpdateUsage(ctx *gin.Context) {
	var form f.UpdateUsageForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	usage := middleware.GetUsage(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.UpdateUsage(tx, usage, form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func RemoveUsage(ctx *gin.Context) {
	usage := middleware.GetUsage(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.DeleteUsage(tx, usage)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
