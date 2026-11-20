package dto

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

type DeleteFileForm struct {
	FileIDs []string `form:"file_ids" json:"file_ids" binding:"required,dive,uuid"`
}

type GetFilesForm struct {
	ListModelsForm
	Type string `form:"type" json:"type" binding:"omitempty,oneof=writeup picture file traffic"`
}

func (f *GetFilesForm) Validate(ctx *gin.Context) model.RetVal {
	if ret := f.ListModelsForm.Validate(ctx); !ret.OK {
		return ret
	}
	return model.SuccessRetVal()
}
