package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	wh "CBCTF/internal/webhook"
	"net/netip"
	"net/url"
	"strings"

	"gorm.io/gorm"
)

func ListWebhooks(tx *gorm.DB, form dto.ListModelsForm) ([]model.Webhook, int64, model.RetVal) {
	return db.InitWebhookRepo(tx).List(form.Limit, form.Offset)
}

func isInWebhookWhitelist(target string) (bool, error) {
	if len(config.Env.Webhook.Whitelist) == 0 {
		return false, nil
	}
	u, err := url.Parse(target)
	if err != nil {
		return false, err
	}
	hostname, err := netip.ParseAddr(u.Hostname())
	if err != nil {
		for _, allowed := range config.Env.Webhook.Whitelist {
			if allowed == u.Hostname() || allowed == u.Host {
				return true, nil
			}
		}
	} else {
		for _, allowed := range config.Env.Webhook.Whitelist {
			if strings.Contains(allowed, "/") {
				prefix, err := netip.ParsePrefix(allowed)
				if err != nil {
					continue
				}
				if prefix.Masked().Contains(hostname) {
					return true, nil
				}
			} else {
				ip, err := netip.ParseAddr(allowed)
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

func validateWebhookURL(target string) model.RetVal {
	in, err := isInWebhookWhitelist(target)
	if err != nil {
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	if !in {
		return model.RetVal{Msg: i18n.Model.Webhook.NotAllowedTarget}
	}
	return model.SuccessRetVal()
}

func CreateWebhook(tx *gorm.DB, form dto.CreateWebhookForm) (model.Webhook, model.RetVal) {
	if ret := validateWebhookURL(form.URL); !ret.OK {
		return model.Webhook{}, ret
	}
	return db.InitWebhookRepo(tx).Create(db.CreateWebhookOptions{
		Name:    form.Name,
		URL:     form.URL,
		Method:  form.Method,
		Headers: form.Headers,
		Timeout: form.Timeout,
		Retry:   form.Retry,
		Events:  form.Events,
	})
}

func UpdateWebhook(tx *gorm.DB, webhook model.Webhook, form dto.UpdateWebhookForm) (model.Webhook, model.RetVal) {
	if form.URL != nil {
		if ret := validateWebhookURL(*form.URL); !ret.OK {
			return model.Webhook{}, ret
		}
	}
	if ret := db.InitWebhookRepo(tx).Update(webhook.ID, db.UpdateWebhookOptions{
		Name:    form.Name,
		URL:     form.URL,
		Method:  form.Method,
		Headers: form.Headers,
		Timeout: form.Timeout,
		Retry:   form.Retry,
		On:      form.On,
		Events:  form.Events,
	}); !ret.OK {
		return model.Webhook{}, ret
	}
	newWebhook, ret := db.InitWebhookRepo(tx).GetByID(webhook.ID)
	if !ret.OK {
		return model.Webhook{}, ret
	}
	wh.DelWebhook(webhook)
	if newWebhook.On {
		wh.AddWebhook(newWebhook)
	}
	return newWebhook, model.SuccessRetVal()
}

func DeleteWebhook(tx *gorm.DB, webhook model.Webhook) model.RetVal {
	if ret := db.InitWebhookRepo(tx).Delete(webhook.ID); !ret.OK {
		return ret
	}
	wh.DelWebhook(webhook)
	return model.SuccessRetVal()
}
