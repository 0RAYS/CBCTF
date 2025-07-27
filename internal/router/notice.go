package router

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	db "CBCTF/internal/repo"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/websocket"
	"CBCTF/internal/websocket/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetNotices(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	DB := db.DB.WithContext(ctx)
	notices, count, ok, msg := db.InitNoticeRepo(DB).List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
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
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
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
	go websocket.SendToAll(false, model.InfoLevel, model.NoticeType, notice.Title, fmt.Sprintf("%d", notice.ContestID))
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": &notice})
}

func UpdateNotice(ctx *gin.Context) {
	var form f.UpdateNoticeForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
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
