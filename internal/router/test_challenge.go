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
			websocket.Send(true, selfID, wm.ErrorLevel, wm.StartVictimWSType, "Start Victim", "Failed")
			victim, ok, _ := db.InitVictimRepo(db.DB).HasAliveVictim(0, challenge.ID)
			if !ok {
				return
			}
			tx := db.DB.Begin()
			if ok, _ = service.StopVictim(tx, victim); !ok {
				tx.Rollback()
				return
			}
			tx.Commit()
			return
		}
		websocket.Send(true, selfID, wm.SuccessLevel, wm.StartVictimWSType, "Start Victim", "Done")
		return
	}()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}

func StopTestVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StopVictimEventType)
	challenge := middleware.GetChallenge(ctx)
	victim, ok, msg := db.InitVictimRepo(db.DB).HasAliveVictim(0, challenge.ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx := db.DB.Begin()
	if ok, msg = service.StopVictim(tx, victim); !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
