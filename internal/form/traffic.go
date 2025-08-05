package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"github.com/gin-gonic/gin"
	"time"
)

type GetTrafficForm struct {
	Start time.Time `form:"start" json:"start" binding:"required"`
	End   time.Time `form:"end" json:"end" binding:"required"`
}

func (f *GetTrafficForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}
