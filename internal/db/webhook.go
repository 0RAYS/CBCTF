package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
)

type WebhookRepo struct {
	BasicRepo[model.Webhook]
}

type CreateWebhookOptions struct {
	Name       string
	URL        string
	Method     string
	Headers    model.StringMap
	Timeout    int64
	RetryCount int64
	RetryDelay int64
	On         bool
	Events     model.StringList
}

func (c CreateWebhookOptions) Convert2Model() model.Model {
	return model.Webhook{
		Name:       c.Name,
		URL:        c.URL,
		Method:     c.Method,
		Headers:    c.Headers,
		Timeout:    c.Timeout,
		RetryCount: c.RetryCount,
		RetryDelay: c.RetryDelay,
		On:         c.On,
		Events:     c.Events,
	}
}

type UpdateWebhookOptions struct {
	Name        *string
	URL         *string
	Method      *string
	Headers     *model.StringMap
	Timeout     *int64
	RetryCount  *int64
	RetryDelay  *int64
	On          *bool
	Events      *model.StringList
	Success     *int64
	SuccessLast *time.Time
	Failure     *int64
	FailureLast *time.Time
}

func (u UpdateWebhookOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Name != nil {
		options["name"] = *u.Name
	}
	if u.URL != nil {
		options["url"] = *u.URL
	}
	if u.Method != nil {
		options["method"] = *u.Method
	}
	if u.Headers != nil {
		options["headers"] = *u.Headers
	}
	if u.Timeout != nil {
		options["timeout"] = *u.Timeout
	}
	if u.RetryCount != nil {
		options["retry_count"] = *u.RetryCount
	}
	if u.RetryDelay != nil {
		options["retry_delay"] = *u.RetryDelay
	}
	if u.On != nil {
		options["on"] = *u.On
	}
	if u.Events != nil {
		options["events"] = *u.Events
	}
	if u.Success != nil {
		options["success"] = *u.Success
	}
	if u.SuccessLast != nil {
		options["success_last"] = *u.SuccessLast
	}
	if u.Failure != nil {
		options["failure"] = *u.Failure
	}
	if u.FailureLast != nil {
		options["failure_last"] = *u.FailureLast
	}
	return options
}

func InitWebhookRepo(tx *gorm.DB) *WebhookRepo {
	return &WebhookRepo{
		BasicRepo: BasicRepo[model.Webhook]{
			DB: tx,
		},
	}
}

func (w *WebhookRepo) Delete(idL ...uint) (bool, string) {
	webhookL, _, ok, msg := w.List(-1, -1, GetOptions{
		Conditions: map[string]any{"id": idL},
		Selects:    []string{"id"},
		Preloads:   map[string]GetOptions{"WebhookHistories": {Selects: []string{"id", "webhook_id"}}},
	})
	if !ok && msg != i18n.WebhookNotFound {
		return false, msg
	}
	webhookHistoryIDL := make([]uint, 0)
	for _, webhook := range webhookL {
		for _, history := range webhook.WebhookHistories {
			webhookHistoryIDL = append(webhookHistoryIDL, history.ID)
		}
	}
	if ok, msg = InitWebhookHistoryRepo(w.DB).Delete(webhookHistoryIDL...); !ok {
		return false, msg
	}
	if res := w.DB.Model(&model.Webhook{}).Where("id IN ?", idL).Delete(&model.Webhook{}); res.Error != nil {
		log.Logger.Warningf("Failed to deleted Webhook: %s", res.Error)
		return false, i18n.DeleteWebhookError
	}
	return true, i18n.Success
}
