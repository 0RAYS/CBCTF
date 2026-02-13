package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCheats(ctx *gin.Context) {
	var form dto.GetCheatsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	options := db.GetOptions{Conditions: map[string]any{"contest_id": middleware.GetContest(ctx).ID}}
	if form.Type != "" {
		options.Conditions["type"] = form.Type
	}
	if form.ReasonType != "" {
		options.Conditions["reason_type"] = form.ReasonType
	}
	cheats, count, ret := db.InitCheatRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	countOptions := db.CountOptions{
		Conditions: options.Conditions,
	}
	countOptions.Conditions["checked"] = true
	checked, ret := db.InitCheatRepo(db.DB).Count(countOptions)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, cheat := range cheats {
		data = append(data, resp.GetCheatResp(cheat))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"count": count, "checked": checked, "cheats": data}))
}

func UpdateCheat(ctx *gin.Context) {
	var form dto.UpdateCheatForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateCheatEventType)
	cheat := middleware.GetCheat(ctx)
	ret := db.InitCheatRepo(db.DB).Update(cheat.ID, db.UpdateCheatRepo{
		Reason:  form.Reason,
		Type:    form.Type,
		Checked: form.Checked,
		Comment: form.Comment,
	})
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}

func DeleteCheat(all bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ret model.RetVal
		if all {
			ctx.Set(middleware.CTXEventTypeKey, model.DeleteAllCheatEventType)
			ret = db.InitCheatRepo(db.DB).DeleteByContestID(middleware.GetContest(ctx).ID)
		} else {
			ctx.Set(middleware.CTXEventTypeKey, model.DeleteCheatEventType)
			cheat := middleware.GetCheat(ctx)
			ret = db.InitCheatRepo(db.DB).Delete(cheat.ID)
		}
		if ret.OK {
			ctx.Set(middleware.CTXEventSuccessKey, true)
		}
		ctx.JSON(http.StatusOK, ret)
	}
}

func CheckCheat(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.ManualCheckCheatEventType)
	contest := middleware.GetContest(ctx)
	service.CheckWebReqIP(db.DB, contest)
	service.CheckVictimReqIP(db.DB, contest)
	service.CheckWrongFlag(db.DB, contest)
	service.CheckSameDevice(db.DB, contest)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}
