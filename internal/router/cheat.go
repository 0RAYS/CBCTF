package router

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
)

func GetCheats(ctx *gin.Context) {
	var form dto.GetCheatsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	cheats, count, checked, ret := service.ListCheats(db.DB.WithContext(ctx.Request.Context()), middleware.GetContest(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, cheat := range cheats {
		data = append(data, resp.GetCheatResp(cheat))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "checked": checked, "cheats": data}))
}

func UpdateCheat(ctx *gin.Context) {
	var form dto.UpdateCheatForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateCheatEventType)
	cheat := middleware.GetCheat(ctx)
	ret := service.UpdateCheat(db.DB.WithContext(ctx.Request.Context()), cheat, form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func DeleteCheat(all bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ret model.RetVal
		if all {
			ctx.Set(middleware.CTXEventTypeKey, model.DeleteAllCheatEventType)
			ret = service.DeleteContestCheats(db.DB.WithContext(ctx.Request.Context()), middleware.GetContest(ctx))
		} else {
			ctx.Set(middleware.CTXEventTypeKey, model.DeleteCheatEventType)
			ret = service.DeleteCheat(db.DB.WithContext(ctx.Request.Context()), middleware.GetCheat(ctx))
		}
		if ret.OK {
			ctx.Set(middleware.CTXEventSuccessKey, true)
		}
		resp.JSON(ctx, ret)
	}
}

func CheckCheat(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.ManualCheckCheatEventType)
	contest := middleware.GetContest(ctx)
	tx := db.DB.WithContext(ctx.Request.Context())
	service.CheckWebReqIP(tx, contest)
	service.CheckVictimReqIP(tx, contest)
	service.CheckWrongFlag(tx, contest)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
