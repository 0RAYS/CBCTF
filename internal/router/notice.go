package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/websocket"
	wsm "CBCTF/internal/websocket/model"
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
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
	ctx.Set(middleware.CTXEventTypeKey, model.CreateNoticeEventType)
	contest := middleware.GetContest(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	notice, ok, msg := service.CreateNotice(tx, contest, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	go func() {
		contestUserIDL := make([]uint, 0)
		for _, user := range contest.Users {
			contestUserIDL = append(contestUserIDL, user.ID)
		}
		idL := make([]uint, 0)
		websocket.UserClientsMu.Lock()
		for id, _ := range websocket.UserClients {
			if slices.Contains(contestUserIDL, id) {
				idL = append(idL, id)
			}
		}
		websocket.UserClientsMu.Unlock()
		websocket.SendToClients(false, wsm.NoticeLevel, wsm.ContestNoticeWSType, fmt.Sprintf("Notice: %s", notice.Title), notice.Content, idL...)
	}()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": &notice})
}

func UpdateNotice(ctx *gin.Context) {
	var form f.UpdateNoticeForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateNoticeEventType)
	notice := middleware.GetNotice(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.UpdateNotice(tx, notice, form)
	if !ok {
		tx.Rollback()
	} else {
		ctx.Set(middleware.CTXEventSuccessKey, true)
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteNotice(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteNoticeEventType)
	notice := middleware.GetNotice(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.DeleteNotice(tx, notice)
	if !ok {
		tx.Rollback()
	} else {
		ctx.Set(middleware.CTXEventSuccessKey, true)
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
