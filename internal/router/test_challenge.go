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
		"remote":   service.GetVictimStatus(db.DB, 0, challenge),
		"file": func() string {
			path := challenge.AttachmentPath(0)
			record, _, _ := db.InitFileRepo(db.DB).Get(db.GetOptions{
				Conditions: map[string]any{"challenge_id": challenge.ID, "type": model.ChallengeFileType}},
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
	go func() {
		_, ok, _ := service.StartVictim(db.DB, 0, 0, 0, challenge.ID)
		if !ok {
			go func() {
				victim, ok, _ := db.InitVictimRepo(db.DB).HasAliveVictim(0, challenge.ID)
				if !ok {
					return
				}
				service.StopVictim(db.DB, victim)
			}()
			websocket.Send(true, selfID, wm.ErrorLevel, wm.StartVictimWSType, "Start Victim", "Failed")
			return
		}
		websocket.Send(true, selfID, wm.SuccessLevel, wm.StartVictimWSType, "Start Victim", "Done")
		return
	}()
	status := service.GetVictimStatus(db.DB, 0, challenge)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": status})
}

func StopTestVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StopVictimEventType)
	challenge := middleware.GetChallenge(ctx)
	victim, ok, msg := db.InitVictimRepo(db.DB).HasAliveVictim(0, challenge.ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ok, msg = service.StopVictim(db.DB, victim)
	if ok {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
