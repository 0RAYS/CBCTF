package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddUsage(ctx *gin.Context) {
	var form f.CreateUsageForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	usages, ok, msg := db.CreateUsage(tx, form, middleware.GetContest(ctx).ID)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": &usages})
}

func GetUsages(ctx *gin.Context) {
	var (
		usages  []model.Usage
		ok      bool
		msg     string
		all     = middleware.GetRole(ctx) == "admin"
		DB      = db.DB.WithContext(ctx)
		contest = middleware.GetContest(ctx)
		team    = middleware.GetTeam(ctx)
	)
	usages, ok, msg = db.GetUsageByContestID(DB, contest.ID, all)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	var challenges []map[string]interface{}
	for _, usage := range usages {
		tmp := map[string]interface{}{}
		challenge, ok, msg := db.GetChallengeByID(DB, usage.ChallengeID)
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
				"solved":   db.IsSolved(DB, contest, team, challenge),
				"attempts": db.CountAttempts(DB, contest, team, challenge),
				"init": func() bool {
					_, ok, _ = db.GetFlagBy3ID(db.DB.WithContext(ctx), contest.ID, team.ID, challenge.ID)
					return ok
				}(),
			}
		}
		challenges = append(challenges, tmp)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": &challenges})
}

func RemoveUsage(ctx *gin.Context) {
	var DB = db.DB.WithContext(ctx)
	usage, ok, msg := db.GetUsageBy2ID(DB, middleware.GetContest(ctx).ID, middleware.GetChallenge(ctx).ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx := DB.Begin()
	ok, msg = db.DeleteUsage(tx, usage.ID)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateUsage(ctx *gin.Context) {
	var DB = db.DB.WithContext(ctx)
	usage, ok, msg := db.GetUsageBy2ID(DB, middleware.GetContest(ctx).ID, middleware.GetChallenge(ctx).ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	var form f.UpdateUsageForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	data := utils.Form2Map(form)
	tx := DB.Begin()
	ok, msg = db.UpdateUsage(tx, usage.ID, data)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
