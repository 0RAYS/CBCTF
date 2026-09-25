package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

func GetNoticeResp(notice model.Notice) gin.H {
	return gin.H{
		"id":         notice.ID,
		"contest_id": notice.ContestID,
		"title":      notice.Title,
		"content":    notice.Content,
		"type":       notice.Type,
		"created_at": notice.UpdatedAt,
	}
}
