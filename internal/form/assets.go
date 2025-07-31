package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"github.com/gin-gonic/gin"
)

type GetAssetForm struct {
	Filename string `form:"filename" json:"filename" binding:"required"`
}

func (f *GetAssetForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}
