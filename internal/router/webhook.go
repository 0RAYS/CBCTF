package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	wh "CBCTF/internal/webhook"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetWebhook(ctx *gin.Context) {
	webhook := middleware.GetWebhook(ctx)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetWebhookResp(webhook)))
}

func GetWebhooks(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	webhooks, count, ret := db.InitWebhookRepo(db.DB).List(form.Limit, form.Offset)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, webhook := range webhooks {
		data = append(data, resp.GetWebhookResp(webhook))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"count": count, "webhooks": data}))
}

func CreateWebhook(ctx *gin.Context) {
	var form dto.CreateWebhookForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateWebhookEventType)
	webhook, ret := db.InitWebhookRepo(db.DB).Create(db.CreateWebhookOptions{
		Name:    form.Name,
		URL:     form.URL,
		Method:  form.Method,
		Headers: form.Headers,
		Timeout: form.Timeout,
		Retry:   form.Retry,
		Events:  form.Events,
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetWebhookResp(webhook)))
}

func UpdateWebhook(ctx *gin.Context) {
	var form dto.UpdateWebhookForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateWebhookEventType)
	webhook := middleware.GetWebhook(ctx)
	if ret := db.InitWebhookRepo(db.DB).Update(webhook.ID, db.UpdateWebhookOptions{
		Name:    form.Name,
		URL:     form.URL,
		Method:  form.Method,
		Headers: form.Headers,
		Timeout: form.Timeout,
		Retry:   form.Retry,
		On:      form.On,
		Events:  form.Events,
	}); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	newWebhook, ret := db.InitWebhookRepo(db.DB).GetByID(webhook.ID)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	wh.DelWebhook(webhook)
	if newWebhook.On {
		wh.AddWebhook(newWebhook)
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}

func DeleteWebhook(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteWebhookEventType)
	webhook := middleware.GetWebhook(ctx)
	if ret := db.InitWebhookRepo(db.DB).Delete(webhook.ID); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	wh.DelWebhook(webhook)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}
