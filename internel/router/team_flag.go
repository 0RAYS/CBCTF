package router

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/k8s"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitTeamFlag(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	teamFlags, ok, msg := service.CreateTeamFlag(tx, team, contestChallenge)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	switch contestChallenge.Challenge.Type {
	case model.DynamicChallengeType:
		ok, msg = k8s.GenerateAttachment(contestChallenge, team, teamFlags)
	default:
		ok, msg = true, i18n.Success
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func ResetTeamFlag(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	teamFlags, ok, msg := service.UpdateTeamFlag(tx, team, contestChallenge)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	switch contestChallenge.Challenge.Type {
	case model.DynamicChallengeType:
		ok, msg = k8s.GenerateAttachment(contestChallenge, team, teamFlags)
	case model.PodsChallengeType:
		// 不考虑失败
		go service.StopVictim(db.DB.WithContext(ctx.Copy()), team, contestChallenge)
		ok, msg = true, i18n.Success
	default:
		ok, msg = true, i18n.Success
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
