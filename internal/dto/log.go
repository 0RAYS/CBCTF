package dto

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type GetLogsForm struct {
	Offset int `form:"offset" json:"offset" binding:"gte=0"`
	Limit  int `form:"limit" json:"limit" binding:"gte=0,lte=100"`
}

func (f *GetLogsForm) Validate(ctx *gin.Context) model.RetVal {
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 100
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	return model.SuccessRetVal()
}
