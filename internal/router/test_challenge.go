package router

import (
	"CBCTF/internal/db"
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
			record, _ := db.InitFileRepo(db.DB).Get(db.GetOptions{
				Conditions: map[string]any{"model": challenge.ModelName(), "model_id": challenge.ID, "type": model.ChallengeFileType}},
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
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}

func StartTestVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StartVictimEventType)
	challenge := middleware.GetChallenge(ctx)
	selfID := middleware.GetSelfID(ctx)
	go func() {
		_, ret := service.StartVictim(db.DB, 0, 0, 0, 0, challenge.ID)
		if !ret.OK {
			websocket.Send(true, selfID, wm.ErrorLevel, wm.StartVictimWSType, "Start Victim", "Failed")
			victim, ret := db.InitVictimRepo(db.DB).HasAliveVictim(0, challenge.ID)
			if !ret.OK {
				return
			}
			service.StopVictim(db.DB, victim)
			return
		}
		websocket.Send(true, selfID, wm.SuccessLevel, wm.StartVictimWSType, "Start Victim", "Done")
		return
	}()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}

func StopTestVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StopVictimEventType)
	challenge := middleware.GetChallenge(ctx)
	victim, ret := db.InitVictimRepo(db.DB).HasAliveVictim(0, challenge.ID)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	if ret = service.StopVictim(db.DB, victim); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, ret)
}
