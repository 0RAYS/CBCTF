package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

// DeleteFileForm for delete files
type DeleteFileForm struct {
	FileIDs []string `form:"file_ids" json:"file_ids" binding:"required,dive,uuid"`
}

func (f *DeleteFileForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

type GetFilesForm struct {
	Offset int    `form:"offset" json:"offset" binding:"gte=0"`
	Limit  int    `form:"limit" json:"limit" binding:"gte=0,lte=100"`
	Type   string `form:"type" json:"type" binding:"omitempty,oneof=writeup picture file traffic"`
}

func (f *GetFilesForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	return model.SuccessRetVal()
}
