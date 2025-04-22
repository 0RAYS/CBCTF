package resp

import (
	"CBCTF/internel/config"
	"CBCTF/internel/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

// GetNoticeResp model.Notice 需要预加载 model.Admin
func GetNoticeResp(notice model.Notice) gin.H {
	return gin.H{
		"id":         notice.ID,
		"title":      notice.Title,
		"content":    notice.Content,
		"created_at": notice.UpdatedAt,
		"author": gin.H{
			"name":   notice.Admin.Name,
			"avatar": fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(notice.Admin.Avatar, "/")),
		},
	}
}
