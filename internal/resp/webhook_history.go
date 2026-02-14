package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetWebhookHistoryResp(history model.WebhookHistory) gin.H {
	return gin.H{
		"id":         history.ID,
		"webhook_id": history.WebhookID,
		"webhook":    history.Webhook.Name,
		"event_id":   history.EventID,
		"event":      history.Event.Type,
		"resp":       history.RespCode,
		"duration":   history.Duration.Milliseconds(),
		"success":    history.Success,
		"error":      history.Error,
		"time":       history.CreatedAt,
	}
}
