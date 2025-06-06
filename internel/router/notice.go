package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetNotices(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 5
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	contest := middleware.GetContest(ctx)
	DB := db.DB.WithContext(ctx)
	notices, count, ok, msg := db.InitNoticeRepo(DB).ListWithConditions(form.Limit, form.Offset, db.GetOptions{
		{
			Key:   "contest_id",
			Value: contest.ID,
			Op:    "and",
		},
	})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, notice := range notices {
		data = append(data, resp.GetNoticeResp(notice))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "notices": data}})
}

func GetNotice(ctx *gin.Context) {
	notice := middleware.GetNotice(ctx)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": resp.GetNoticeResp(notice)})
}

func CreateNotice(ctx *gin.Context) {
	var form f.CreateNoticeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	notice, ok, msg := service.CreateNotice(tx, contest, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": &notice})
}

func UpdateNotice(ctx *gin.Context) {
	var form f.UpdateNoticeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	notice := middleware.GetNotice(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.UpdateNotice(tx, notice, form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteNotice(ctx *gin.Context) {
	notice := middleware.GetNotice(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.DeleteNotice(tx, notice)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
