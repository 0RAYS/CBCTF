package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"

	"github.com/gin-gonic/gin"
)

type GetLogsForm struct {
	Offset int `form:"offset" json:"offset"`
	Limit  int `form:"limit" json:"limit"`
}

func (f *GetLogsForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 100
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	return true, i18n.Success
}
