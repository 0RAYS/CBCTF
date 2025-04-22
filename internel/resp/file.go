package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetFileResp(file model.File) gin.H {
	return gin.H{
		"id":       file.ID,
		"filename": file.Filename,
		"hash":     file.Hash,
		"size":     file.Size,
		"date":     file.CreatedAt,
		"suffix":   file.Suffix,
		"uploader": file.Uploader,
	}
}
