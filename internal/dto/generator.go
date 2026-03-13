package dto

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type ListGeneratorsForm struct {
	ListModelsForm
	Deleted bool `form:"deleted" json:"deleted"`
}

func (f *ListGeneratorsForm) Validate(ctx *gin.Context) model.RetVal {
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	return model.SuccessRetVal()
}

type StartGeneratorsForm struct {
	Challenges []string `form:"challenges" json:"challenges" binding:"required,dive,uuid"`
}

type StopGeneratorsForm struct {
	Generators []uint `form:"generators" json:"generators" binding:"required,dive,gt=0"`
}
