package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	wh "CBCTF/internal/webhook"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetWebhook(ctx *gin.Context) {
	webhook := middleware.GetWebhook(ctx)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": resp.GetWebhookResp(webhook)})
}

func GetWebhooks(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	webhooks, count, ok, msg := db.InitWebhookRepo(db.DB).List(form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, webhook := range webhooks {
		data = append(data, resp.GetWebhookResp(webhook))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": gin.H{"count": count, "webhooks": data}})
}

func CreateWebhook(ctx *gin.Context) {
	var form f.CreateWebhookForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateWebhookEventType)
	webhook, ok, msg := db.InitWebhookRepo(db.DB).Create(db.CreateWebhookOptions{
		Name:    form.Name,
		URL:     form.URL,
		Method:  form.Method,
		Headers: form.Headers,
		Timeout: form.Timeout,
		Retry:   form.Retry,
		Events:  form.Events,
	})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": resp.GetWebhookResp(webhook)})
}

func UpdateWebhook(ctx *gin.Context) {
	var form f.UpdateWebhookForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateWebhookEventType)
	webhook := middleware.GetWebhook(ctx)
	if ok, msg := db.InitWebhookRepo(db.DB).Update(webhook.ID, db.UpdateWebhookOptions{
		Name:    form.Name,
		URL:     form.URL,
		Method:  form.Method,
		Headers: form.Headers,
		Timeout: form.Timeout,
		Retry:   form.Retry,
		On:      form.On,
		Events:  form.Events,
	}); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	newWebhook, ok, msg := db.InitWebhookRepo(db.DB).GetByID(webhook.ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	wh.DelWebhook(webhook)
	if newWebhook.On {
		wh.AddWebhook(newWebhook)
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}

func DeleteWebhook(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteWebhookEventType)
	webhook := middleware.GetWebhook(ctx)
	if ok, msg := db.InitWebhookRepo(db.DB).Delete(webhook.ID); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	wh.DelWebhook(webhook)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}
