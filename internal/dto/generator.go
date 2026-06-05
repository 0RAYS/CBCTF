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

type GetGeneratorLogsForm struct {
	Lines int64 `form:"lines" json:"lines" binding:"omitempty,gt=0"`
}

func (f *GetGeneratorLogsForm) Validate(_ *gin.Context) model.RetVal {
	if f.Lines <= 0 {
		f.Lines = 1000
	}
	return model.SuccessRetVal()
}
