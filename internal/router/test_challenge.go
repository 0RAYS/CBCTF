package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"CBCTF/internal/websocket"
	wm "CBCTF/internal/websocket/model"
	"net/http"
	"os"

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
			path := challenge.AttachmentPath(0)
			record, _, _ := db.InitFileRepo(db.DB.WithContext(ctx)).Get(db.GetOptions{
				Conditions: map[string]any{"challenge_id": challenge.ID, "type": model.ChallengeFile}},
			)
			filename := "attachment.zip"
			if record.Path == path {
				filename = record.Filename
			}
			if _, err := os.Stat(path); err != nil {
				return ""
			}
			return filename
		}(),
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
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
