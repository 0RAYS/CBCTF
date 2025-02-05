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
	tx := db.DB.WithContext(ctx).Begin()
	usages, ok, msg := db.CreateUsage(tx, ctx, form, middleware.GetContest(ctx).ID)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": usages})
}

func GetUsages(ctx *gin.Context) {
	var (
		usages []model.Usage
		ok     bool
		msg    string
		all    = middleware.GetRole(ctx) == "admin"
	)
	usages, ok, msg = db.GetUsageByContestID(ctx, middleware.GetContest(ctx).ID, all)
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
		if !all {
			tmp["status"] = gin.H{
				"solved":   db.IsSolved(ctx, middleware.GetContest(ctx), middleware.GetTeam(ctx), challenge),
				"attempts": db.CountAttempts(ctx, middleware.GetContest(ctx), middleware.GetTeam(ctx), challenge),
			}
		}
		challenges = append(challenges, tmp)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": challenges})
}

func RemoveUsage(ctx *gin.Context) {
	usage, ok, msg := db.GetUsageBy2ID(ctx, middleware.GetContest(ctx).ID, middleware.GetChallenge(ctx).ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg = db.DeleteUsage(tx, usage.ID)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateUsage(ctx *gin.Context) {
	usage, ok, msg := db.GetUsageBy2ID(ctx, middleware.GetContest(ctx).ID, middleware.GetChallenge(ctx).ID)
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
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg = db.UpdateUsage(tx, usage.ID, data)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
