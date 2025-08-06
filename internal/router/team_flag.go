package router

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	db "CBCTF/internal/repo"
	"CBCTF/internal/service"
	"CBCTF/internal/websocket"
	wm "CBCTF/internal/websocket/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitTeamFlag(ctx *gin.Context) {
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
	switch contestChallenge.Type {
	case model.DynamicChallengeType:
		go func() {
			if ok, _ = k8s.GenerateAttachment(contestChallenge, team, teamFlags); !ok {
				websocket.Send(false, user.ID, wm.ErrorLevel, wm.GenerateAttachmentType, "Generate Attachment", "Failed")
				return
			}
			websocket.Send(false, user.ID, wm.SuccessLevel, wm.GenerateAttachmentType, "Generate Attachment", "Done")
		}()
		ok, msg = true, i18n.Success
	default:
		ok, msg = true, i18n.Success
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func ResetTeamFlag(ctx *gin.Context) {
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
	tx.Commit()
	switch contestChallenge.Type {
	case model.DynamicChallengeType:
		go func() {
			if ok, _ = k8s.GenerateAttachment(contestChallenge, team, teamFlags); !ok {
				websocket.Send(false, user.ID, wm.ErrorLevel, wm.GenerateAttachmentType, "Generate Attachment", "Failed")
				return
			}
			websocket.Send(false, user.ID, wm.SuccessLevel, wm.GenerateAttachmentType, "Generate Attachment", "Done")
		}()
		ok, msg = true, i18n.Success
	case model.PodsChallengeType:
		// 不考虑失败
		go service.StopTeamVictim(db.DB.WithContext(ctx.Copy()), team, contestChallenge)
		ok, msg = true, i18n.Success
	default:
		ok, msg = true, i18n.Success
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
