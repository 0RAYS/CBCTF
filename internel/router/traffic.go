package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTraffics(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 10
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	victim := middleware.GetVictim(ctx)
	repo := db.InitTrafficRepo(db.DB.WithContext(ctx))
	traffics, count, ok, msg := repo.GetByKey("victim_id", victim.ID, form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, traffic := range traffics {
		data = append(data, resp.GetTrafficResp(traffic))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"traffics": data, "count": count}})
}
