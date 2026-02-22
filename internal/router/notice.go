package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/websocket"
	wsm "CBCTF/internal/websocket/model"
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func GetNotices(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	contest := middleware.GetContest(ctx)
	notices, count, ret := db.InitNoticeRepo(db.DB).List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, notice := range notices {
		data = append(data, resp.GetNoticeResp(notice))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"count": count, "notices": data}))
}

func CreateNotice(ctx *gin.Context) {
	var form dto.CreateNoticeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateNoticeEventType)
	contest := middleware.GetContest(ctx)
	notice, ret := db.InitNoticeRepo(db.DB).Create(db.CreateNoticeOptions{
		ContestID: contest.ID,
		Title:     form.Title,
		Content:   form.Content,
		Type:      form.Type,
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	go func() {
		contestUserIDL, ret := db.InitUserRepo(db.DB).GetIDByContestID(contest.ID, -1, -1)
		if !ret.OK {
			return
		}
		idL := make([]uint, 0)
		websocket.UserClientsMu.Lock()
		for id := range websocket.UserClients {
			if slices.Contains(contestUserIDL, id) {
				idL = append(idL, id)
			}
		}
		websocket.UserClientsMu.Unlock()
		websocket.SendToClients(false, wsm.NoticeLevel, wsm.ContestNoticeWSType, fmt.Sprintf("Notice: %s", notice.Title), notice.Content, idL...)
	}()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(notice))
}

func UpdateNotice(ctx *gin.Context) {
	var form dto.UpdateNoticeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateNoticeEventType)
	notice := middleware.GetNotice(ctx)
	ret := db.InitNoticeRepo(db.DB).Update(notice.ID, db.UpdateNoticeOptions{
		Title:   form.Title,
		Content: form.Content,
		Type:    form.Type,
	})
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}

func DeleteNotice(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteNoticeEventType)
	notice := middleware.GetNotice(ctx)
	ret := db.InitNoticeRepo(db.DB).Delete(notice.ID)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}
