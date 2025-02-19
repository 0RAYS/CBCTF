package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/middleware"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateNotice(ctx *gin.Context) {
	var form f.CreateNoticeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	notice, ok, msg := db.CreateNotice(tx, contest, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": &notice})
}

func GetNotice(ctx *gin.Context) {
	notice := middleware.GetNotice(ctx)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": &notice})
}

func GetNotices(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	notices, count, ok, msg := db.GetNotices(db.DB.WithContext(ctx), form.Limit, form.Offset, contest.ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "notices": &notices}})
}

func UpdateNotice(ctx *gin.Context) {
	var form f.UpdateNoticeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	notice := middleware.GetNotice(ctx)
	data := utils.Form2Map(form)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.UpdateNotice(tx, notice.ID, data)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}

func DeleteNotice(ctx *gin.Context) {
	notice := middleware.GetNotice(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.DeleteNotice(tx, notice.ID)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}
