package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetNoticesResp(notices []model.Notice) []gin.H {
	data := make([]gin.H, 0)
	for _, notice := range notices {
		data = append(data, gin.H{
			"title":   notice.Title,
			"content": notice.Content,
			"created": notice.CreatedAt,
		})
	}
	return data
}
