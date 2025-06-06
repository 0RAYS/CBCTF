package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetNoticeResp(notice model.Notice) gin.H {
	return gin.H{
		"id":         notice.ID,
		"title":      notice.Title,
		"content":    notice.Content,
		"type":       notice.Type,
		"created_at": notice.UpdatedAt,
	}
}
