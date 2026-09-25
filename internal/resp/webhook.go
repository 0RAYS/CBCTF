package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

func GetWebhookResp(webhook model.Webhook) gin.H {
	return gin.H{
		"id":           webhook.ID,
		"name":         webhook.Name,
		"url":          webhook.URL,
		"method":       webhook.Method,
		"headers":      webhook.Headers,
		"timeout":      webhook.Timeout,
		"retry":        webhook.Retry,
		"on":           webhook.On,
		"events":       webhook.Events,
		"success":      webhook.Success,
		"success_last": webhook.SuccessLast,
		"failure":      webhook.Failure,
		"failure_last": webhook.FailureLast,
	}
}
