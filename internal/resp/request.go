package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetRequestResp(request model.Request) gin.H {
	return gin.H{
		"id":         request.ID,
		"ip":         request.IP,
		"time":       request.Time,
		"method":     request.Method,
		"path":       request.Path,
		"url":        request.URL,
		"user_agent": request.UserAgent,
		"status":     request.Status,
		"referer":    request.Referer,
		"magic":      request.Magic,
		"user_id":    request.UserID.V,
	}
}
