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
	usages, _, msg := db.CreateUsage(ctx, form, middleware.GetContestID(ctx))
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": usages})
}

func GetUsages(ctx *gin.Context) {
	var (
		usages []model.Usage
		ok     bool
		msg    string
		all    = middleware.GetRole(ctx) == "admin"
	)
	usages, ok, msg = db.GetUsageByContestID(ctx, middleware.GetContestID(ctx), all)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	var challenges []map[string]interface{}
	for _, usage := range usages {
		tmp := map[string]interface{}{}
		challenge, ok, msg := db.GetChallengeByID(ctx, usage.ChallengeID)
		if !ok {
			log.Logger.Warningf("Failed to get challenge %s: %s", usage.ChallengeID, msg)
			continue
		}
		if !all {
			usage.Flag = ""
			challenge.Flag = ""
			challenge.DockerImage = ""
			challenge.GeneratorImage = ""
		}
		tmp["usage"] = usage
		tmp["challenge"] = challenge
		challenges = append(challenges, tmp)
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
