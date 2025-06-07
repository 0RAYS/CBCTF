package router

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/middleware"
	"CBCTF/internel/redis"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitTeamFlag(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	if ok, err := redis.CheckChallengeInit(team.ID, contestChallenge.ID); !ok || err != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"msg": i18n.TooManyRequests, "data": nil})
		return
	}
	_ = redis.RecordChallengeInit(team.ID, contestChallenge.ID)
	tx := db.DB.WithContext(ctx).Begin()
	_, ok, msg := service.CreateTeamFlag(tx, team, contestChallenge)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func ResetTeamFlag(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	if ok, err := redis.CheckChallengeInit(team.ID, contestChallenge.ID); !ok || err != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"msg": i18n.TooManyRequests, "data": nil})
		return
	}
	_ = redis.RecordChallengeInit(team.ID, contestChallenge.ID)
	tx := db.DB.WithContext(ctx).Begin()
	_, ok, msg := service.UpdateTeamFlag(tx, team, contestChallenge)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	// TODO
	//switch usage.Challenge.Type {
	//case model.DynamicChallenge:
	//	ok, msg = k8s.GenerateAttachment(usage, team, answers)
	//case model.PodsChallenge:
	//	// 不考虑失败
	//	go service.StopVictim(db.DB.WithContext(ctx.Copy()), team, usage)
	//	ok, msg = true, i18n.Success
	//default:
	//	ok, msg = true, i18n.Success
	//}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
