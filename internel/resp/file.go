package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetFileResp(file model.File) gin.H {
	return gin.H{
		"filename": file.Filename,
		"hash":     file.Hash,
		"size":     file.Size,
		"id":       file.ID,
		"date":     file.CreatedAt,
		"suffix":   file.Suffix,
		"uploader": file.Uploader,
	}
}
