package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"

	"github.com/gin-gonic/gin"
)

func GetTestChallengeStatus(ctx *gin.Context) {
	status := service.GetTestChallengeStatus(db.DB, middleware.GetChallenge(ctx))
	resp.JSON(ctx, model.SuccessRetVal(resp.GetContestChallengeStatusResp(status)))
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
	ret := service.StopTestVictim(db.DB, middleware.GetChallenge(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, ret)
}
