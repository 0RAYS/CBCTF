package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetEventResp(event model.Event) gin.H {
	return gin.H{
		"id":         event.ID,
		"type":       event.Type,
		"success":    event.Success,
		"ip":         event.IP,
		"magic":      event.Magic,
		"models":     event.Models,
		"created_at": event.CreatedAt,
	}
}
