package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCheats(ctx *gin.Context) {
	var form f.GetCheatsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	options := db.GetOptions{
		Conditions: map[string]any{"contest_id": middleware.GetContest(ctx).ID},
	}
	if form.Type != "" {
		options.Conditions["type"] = form.Type
	}
	cheats, count, ok, msg := db.InitCheatRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	countOptions := db.CountOptions{
		Conditions: options.Conditions,
	}
	countOptions.Conditions["checked"] = true
	checked, ok, msg := db.InitCheatRepo(db.DB).Count(countOptions)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, cheat := range cheats {
		data = append(data, resp.GetCheatResp(cheat))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "checked": checked, "cheats": data}})
}

func UpdateCheat(ctx *gin.Context) {
	var form f.UpdateCheatForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateCheatEventType)
	cheat := middleware.GetCheat(ctx)
	ok, msg := db.InitCheatRepo(db.DB).Update(cheat.ID, db.UpdateCheatRepo{
		Reason:  form.Reason,
		Type:    form.Type,
		Checked: form.Checked,
		Comment: form.Comment,
	})
	if ok {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
