package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"CBCTF/internal/task"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitTeamFlag(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.InitChallengeEventType)
	team := middleware.GetTeam(ctx)
	user := middleware.GetSelf(ctx).(model.User)
	contestChallenge := middleware.GetContestChallenge(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	teamFlags, ok, msg := service.CreateTeamFlag(tx, team, contestChallenge)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if contestChallenge.Type == model.DynamicChallengeType {
		if _, err := task.EnqueueGenAttachmentTask(user.ID, contestChallenge, team, teamFlags); err != nil {
			log.Logger.Warningf("Enqueue gen attachment task failed %v", err)
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.EnqueueTaskError, "data": nil})
			return
		}
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}

func ResetTeamFlag(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.ResetChallengeEventType)
	team := middleware.GetTeam(ctx)
	user := middleware.GetSelf(ctx).(model.User)
	contestChallenge := middleware.GetContestChallenge(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	teamFlags, ok, msg := service.UpdateTeamFlag(tx, team, contestChallenge)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	switch contestChallenge.Type {
	case model.DynamicChallengeType:
		if _, err := task.EnqueueGenAttachmentTask(user.ID, contestChallenge, team, teamFlags); err != nil {
			log.Logger.Warningf("Enqueue gen attachment task failed %v", err)
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.EnqueueTaskError, "data": nil})
			return
		}
	case model.PodsChallengeType:
		// 不考虑失败
		go service.StopTeamVictim(db.DB.WithContext(ctx.Copy()), team, contestChallenge)
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
