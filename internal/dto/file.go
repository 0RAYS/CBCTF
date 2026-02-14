package dto

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

// DeleteFileForm for delete files
type DeleteFileForm struct {
	FileIDs []string `form:"file_ids" json:"file_ids" binding:"required,dive,uuid"`
}

type GetFilesForm struct {
	Offset int    `form:"offset" json:"offset" binding:"gte=0"`
	Limit  int    `form:"limit" json:"limit" binding:"gte=0,lte=100"`
	Type   string `form:"type" json:"type" binding:"omitempty,oneof=writeup picture file traffic"`
}

func (f *GetFilesForm) Validate(ctx *gin.Context) model.RetVal {
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	return model.SuccessRetVal()
}
