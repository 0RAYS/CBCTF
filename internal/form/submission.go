package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"github.com/gin-gonic/gin"
)

// SubmitFlagForm for submit flag
type SubmitFlagForm struct {
	Flag string `form:"flag" json:"flag" binding:"required"`
}

func (f *SubmitFlagForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}
