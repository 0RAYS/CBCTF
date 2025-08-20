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
	go func(ctx *gin.Context) {
		websocket.UserClientsMu.Lock()
		userIDL := make([]uint, 0)
		for userID := range websocket.UserClients {
			userIDL = append(userIDL, userID)
		}
		websocket.UserClientsMu.Unlock()
		contest, ok, msg = db.InitContestRepo(db.DB.WithContext(ctx)).GetByID(contest.ID, db.GetOptions{
			Selects:  []string{"id"},
			Preloads: map[string]db.GetOptions{"Users": {Conditions: map[string]any{"id": userIDL}, Selects: []string{"id"}}},
		})
		contestUserIDL := make([]uint, 0)
		for _, user := range contest.Users {
			contestUserIDL = append(contestUserIDL, user.ID)
		}
		websocket.SendToClients(false, wsm.NoticeLevel, wsm.ContestNoticeWSType, fmt.Sprintf("Notice: %s", notice.Title), notice.Content, userIDL...)
	}(ctx.Copy())
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
