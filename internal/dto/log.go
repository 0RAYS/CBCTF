package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type GetLogsForm struct {
	Offset int `form:"offset" json:"offset"`
	Limit  int `form:"limit" json:"limit"`
}

func (f *GetLogsForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	if f.Limit < 0 {
		f.Limit = 100
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 100
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	return model.SuccessRetVal()
}
