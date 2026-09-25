package dto

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

// ListModelsForm for get models list
type ListModelsForm struct {
	Offset int `form:"offset" json:"offset" binding:"gte=0"`
	Limit  int `form:"limit" json:"limit" binding:"gte=0,lte=100"`
}

func (f *ListModelsForm) Validate(ctx *gin.Context) model.RetVal {
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	return model.SuccessRetVal()
}
