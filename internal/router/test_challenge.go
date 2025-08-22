package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"CBCTF/internal/websocket"
	wm "CBCTF/internal/websocket/model"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func GetTestChallengeStatus(ctx *gin.Context) {
	challenge := middleware.GetChallenge(ctx)
	data := gin.H{
		"attempts": 0,
		"init":     true,
		"solved":   false,
		"remote":   service.GetTestVictimStatus(db.DB.WithContext(ctx), challenge),
		"file": func() string {
			if _, err := os.Stat(challenge.AttachmentPath(0)); err != nil {
				return ""
			}
			return "attachment.zip"
		}(),
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
}

func DownloadTestAttachment(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DownloadAttachmentEventType)
	challenge := middleware.GetChallenge(ctx)
	if challenge.Type == model.DynamicChallengeType {
		_ = os.Remove(challenge.AttachmentPath(0))
		if ok, msg := service.GenTestAttachment(db.DB.WithContext(ctx), challenge); !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		for i := 0; i < 10; i++ {
			if _, err := os.Stat(challenge.AttachmentPath(0)); err == nil {
				break
			}
			time.Sleep(time.Second)
		}
	}
	path := challenge.AttachmentPath(0)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.FileNotFound, "data": nil})
			return
		}
		log.Logger.Warningf("Failed to get attachment: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.File(path)
}

func StartTestVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StartVictimEventType)
	challenge := middleware.GetChallenge(ctx)
	selfID := middleware.GetSelfID(ctx)
	go func(ctx *gin.Context) {
		tx := db.DB.WithContext(ctx).Begin()
		_, ok, _ := service.StartTestVictim(tx, challenge)
		if !ok {
			go service.StopTestVictim(db.DB.WithContext(ctx.Copy()), challenge)
			tx.Rollback()
			websocket.Send(true, selfID, wm.ErrorLevel, wm.StartVictimWSType, "Start Victim", "Failed")
			return
		}
		tx.Commit()
		websocket.Send(true, selfID, wm.SuccessLevel, wm.StartVictimWSType, "Start Victim", "Done")
		return
	}(ctx.Copy())
	status := service.GetTestVictimStatus(db.DB.WithContext(ctx), challenge)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": status})
}

func StopTestVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StopVictimEventType)
	challenge := middleware.GetChallenge(ctx)
	_, msg := service.StopTestVictim(db.DB.WithContext(ctx), challenge)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
