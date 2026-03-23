package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
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
				Conditions: map[string]any{"model": model.ModelName(challenge), "model_id": challenge.ID, "type": model.ChallengeFileType}},
			)
			filename := "attachment.zip"
			if string(record.Path) == path {
				filename = record.Filename
			}
			if _, err := os.Stat(path); err != nil {
				return ""
			}
			return filename
		}(),
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func StartTestVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StartVictimEventType)
	challenge := middleware.GetChallenge(ctx)
	selfID := middleware.GetSelf(ctx).ID
	ret := service.StartVictim(db.DB, selfID, 0, 0, 0, challenge.ID)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func StopTestVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StopVictimEventType)
	challenge := middleware.GetChallenge(ctx)
	victim, ret := db.InitVictimRepo(db.DB).HasAliveVictim(0, challenge.ID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if ret = service.StopVictim(db.DB, victim); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, ret)
}
