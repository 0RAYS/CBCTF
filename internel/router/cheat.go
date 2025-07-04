package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCheats(ctx *gin.Context) {
	var form f.GetCheatsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	options := db.GetOptions{}
	if form.Type != "" {
		options.Conditions = map[string]any{"type": form.Type}
	}
	cheats, count, ok, msg := db.InitCheatRepo(db.DB.WithContext(ctx)).List(form.Limit, form.Offset, options)
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
