package router

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddUsage(ctx *gin.Context) {
	var form constants.CreateUsageForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	usage, ok, msg := db.GetUsageBy2ID(ctx, middleware.GetContestID(ctx), form.ChallengeID)
	if ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": "UsageAlreadyExist", "data": nil})
		return
	}
	usage, ok, msg = db.CreateUsage(ctx, form, middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": usage})
}

func GetUsages(ctx *gin.Context) {
	usages, ok, msg := db.GetUsageByContestID(ctx, middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	var challenges []model.Challenge
	for _, usage := range usages {
		challenge, ok, msg := db.GetChallengeByID(ctx, usage.ChallengeID)
		if !ok {
			log.Logger.Warningf("Failed to get challenge %d: %s", usage.ChallengeID, msg)
			continue
		}
		challenges = append(challenges, challenge)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": challenges})
}

func RemoveUsage(ctx *gin.Context) {
	usage, ok, msg := db.GetUsageBy2ID(ctx, middleware.GetContestID(ctx), middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	_, msg = db.DeleteUsage(ctx, usage.ID)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateUsage(ctx *gin.Context) {
	usage, ok, msg := db.GetUsageBy2ID(ctx, middleware.GetContestID(ctx), middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	var form constants.UpdateUsageForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	data := utils.Form2Map(form)
	_, msg = db.UpdateUsage(ctx, usage.ID, data)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
