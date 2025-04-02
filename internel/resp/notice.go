package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetNoticeResp(notice model.Notice) gin.H {
	return gin.H{
		"title":      notice.Title,
		"content":    notice.Content,
		"created_at": notice.CreatedAt,
	}
}
