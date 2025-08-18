package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetWebhookHistoryResp(history model.WebhookHistory) gin.H {
	return gin.H{
		"id":         history.ID,
		"webhook":    history.Webhook.Name,
		"event":      history.Event.Type,
		"resp":       history.RespCode,
		"duration":   history.Duration,
		"success":    history.Success,
		"error":      history.Error,
		"retry":      history.RetryCount,
		"next_retry": history.NextRetry,
		"time":       history.CreatedAt,
	}
}
