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
	user := middleware.GetSelf(ctx).(model.User)
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	contestFlags, _, ok, msg := db.InitContestFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	contestChallenge.ContestFlags = contestFlags
	challenge := middleware.GetChallenge(ctx)
	tx := db.DB.Begin()
	teamFlags, ok, msg := service.CreateTeamFlag(tx, team, contest, contestChallenge)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if challenge.Type == model.DynamicChallengeType {
		if _, err := task.EnqueueGenAttachmentTask(user.ID, challenge, team, teamFlags); err != nil {
			log.Logger.Warningf("Failed to enqueue gen attachment task: %s", err)
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
	user := middleware.GetSelf(ctx).(model.User)
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	contestFlags, _, ok, msg := db.InitContestFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	contestChallenge.ContestFlags = contestFlags
	challenge := middleware.GetChallenge(ctx)
	tx := db.DB.Begin()
	teamFlags, ok, msg := service.UpdateTeamFlag(tx, team, contest, contestChallenge)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	switch challenge.Type {
	case model.DynamicChallengeType:
		if _, err := task.EnqueueGenAttachmentTask(user.ID, challenge, team, teamFlags); err != nil {
			log.Logger.Warningf("Failed to enqueue gen attachment task: %s", err)
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.EnqueueTaskError, "data": nil})
			return
		}
		tx.Commit()
	case model.PodsChallengeType:
		tx.Commit()
		// 不考虑失败
		go func() {
			victim, ok, _ := db.InitVictimRepo(db.DB).HasAliveVictim(team.ID, challenge.ID)
			if !ok {
				return
			}
			service.StopVictim(db.DB, victim)
		}()
	default:
		tx.Commit()
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
