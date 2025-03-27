package router

import (
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
	data["attempts"], _, _ = service.CountAttempts(DB, team, usage)
	data["init"], _, _ = service.IsGenerated(DB, usage, team)
	data["solved"], _, _ = service.IsSolved(DB, team, usage)
	data["remote"] = service.GetRemoteStatus(DB, usage)
	data["file"] = func() string {
		if _, err := os.Stat(usage.Challenge.AttachmentPath(team.ID)); err != nil {
			return ""
		}
		return usage.Challenge.AttachmentPath(team.ID)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}

func GetUsages(ctx *gin.Context) {
	var (
		all     = middleware.GetRole(ctx) == "admin"
		DB      = db.DB.WithContext(ctx)
		contest = middleware.GetContest(ctx)
		team    = middleware.GetTeam(ctx)
	)
	usages, ok, msg := service.GetUsages(DB, contest, all)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	var data []gin.H
	for _, usage := range usages {
		tmp := resp.GetUsageResp(usage)
		tmp["attempts"], _, _ = service.CountAttempts(DB, team, usage)
		tmp["init"], _, _ = service.IsGenerated(DB, usage, team)
		tmp["solved"], _, _ = service.IsSolved(DB, team, usage)
		tmp["remote"] = service.GetRemoteStatus(DB, usage)
		tmp["file"] = func() string {
			if _, err := os.Stat(usage.Challenge.AttachmentPath(team.ID)); err != nil {
				return ""
			}
			return usage.Challenge.AttachmentPath(team.ID)
		}
		data = append(data, tmp)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}

func InitUsage(ctx *gin.Context) {
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
	answers, ok, msg := service.InitAnswer(tx, usage, team)
	if !ok {
		tx.Rollback()
	}
	if usage.Challenge.Type == model.DynamicChallenge {
		ok, msg = service.GeneratorAttachment(usage, answers[0])
	}
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
