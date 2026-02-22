package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	wh "CBCTF/internal/webhook"
	"net/http"
	"net/netip"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetWebhooks(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
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

func checkWebhookBlacklist(target string) (bool, error) {
	u, err := url.Parse(target)
	if err != nil {
		return false, err
	}
	hostname, err := netip.ParseAddr(u.Hostname())
	if err != nil {
		for _, banned := range config.Env.Webhook.Blacklist {
			if banned == u.Hostname() || banned == u.Host {
				return true, nil
			}
		}
	} else {
		for _, banned := range config.Env.Webhook.Blacklist {
			if strings.Contains(banned, "/") {
				prefix, err := netip.ParsePrefix(banned)
				if err != nil {
					continue
				}
				if prefix.Masked().Contains(hostname) {
					return true, nil
				}
			} else {
				ip, err := netip.ParseAddr(banned)
				if err != nil {
					continue
				}
				if ip.Unmap() == hostname {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func CreateWebhook(ctx *gin.Context) {
	var form dto.CreateWebhookForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateWebhookEventType)
	banned, err := checkWebhookBlacklist(form.URL)
	if err != nil {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	if banned {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Model.Webhook.NotAllowedTarget})
		return
	}
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
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateWebhookEventType)
	if form.URL != nil {
		banned, err := checkWebhookBlacklist(*form.URL)
		if err != nil {
			ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
			return
		}
		if banned {
			ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Model.Webhook.NotAllowedTarget})
			return
		}
	}
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
