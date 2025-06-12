package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCheats(ctx *gin.Context) {
	var form f.GetCheatsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 5
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	if form.Type != model.Suspicious && form.Type != model.Cheater {
		form.Type = ""
	}
	conditions := make(db.GetOptions, 0)
	if form.Type != "" {
		conditions = append(conditions, db.GetOption{Key: "type", Value: form.Type, Op: "and"})
	}
	cheats, count, ok, msg := db.InitCheatRepo(db.DB.WithContext(ctx)).ListWithConditions(form.Limit, form.Offset, conditions, false)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, cheat := range cheats {
		data = append(data, resp.GetCheatResp(cheat))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "cheats": data}})
}

func GetCheat(ctx *gin.Context) {
	cheat := middleware.GetCheat(ctx)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": resp.GetCheatResp(cheat)})
}
