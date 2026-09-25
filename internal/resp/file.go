package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

func GetFileResp(file model.File) gin.H {
	return gin.H{
		"id":       file.RandID,
		"filename": file.Filename,
		"type":     file.Type,
		"hash":     file.Hash,
		"size":     file.Size,
		"date":     file.CreatedAt,
		"suffix":   file.Suffix,
		"model":    file.Model,
		"model_id": file.ModelID,
	}
}
