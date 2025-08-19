package db

import (
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
	RetryCount int
	On         bool
	Events     model.StringList
}

func (c CreateWebhookOptions) Convert2Model() model.Model {
	return model.Webhook{
		Name:        c.Name,
		URL:         c.URL,
		Method:      c.Method,
		Headers:     c.Headers,
		Timeout:     c.Timeout,
		RetryCount:  c.RetryCount,
		On:          c.On,
		Events:      c.Events,
		SuccessLast: time.Now(),
		FailureLast: time.Now(),
	}
}

type UpdateWebhookOptions struct {
	Name        *string
	URL         *string
	Method      *string
	Headers     *model.StringMap
	Timeout     *int64
	RetryCount  *int
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

func (w *WebhookRepo) UpdateStatus(id uint, success bool, last time.Time) (bool, string) {
	old, ok, msg := w.GetByID(id, GetOptions{Selects: []string{"id", "success", "failure"}})
	if !ok {
		return false, msg
	}
	var options UpdateWebhookOptions
	if success {
		count := old.Success + 1
		options = UpdateWebhookOptions{
			Success:     &count,
			SuccessLast: &last,
		}
	} else {
		count := old.Failure + 1
		options = UpdateWebhookOptions{
			Failure:     &count,
			FailureLast: &last,
		}
	}
	return w.Update(id, options)
}
