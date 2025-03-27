package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetNotices(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	DB := db.DB.WithContext(ctx)
	notices, count, ok, msg := db.InitNoticeRepo(DB).GetAll(contest.ID, form.Limit, form.Offset, false, 0)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	var data gin.H
	data["notices"] = make([]gin.H, 0)
	for _, notice := range notices {
		data["notices"] = append(data["notices"].([]gin.H), resp.GetNoticeResp(notice))
	}
	data["count"] = count
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": data})
}

func GetNotice(ctx *gin.Context) {
	notice := middleware.GetNotice(ctx)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": resp.GetNoticeResp(notice)})
}
