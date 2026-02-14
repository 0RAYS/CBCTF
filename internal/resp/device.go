package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetDeviceResp(device model.Device) gin.H {
	return gin.H{
		"id":         device.ID,
		"user_id":    device.UserID,
		"magic":      device.Magic,
		"count":      device.Count,
		"created_at": device.CreatedAt,
	}
}
