package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"

	"github.com/gin-gonic/gin"
)

type GetTrafficForm struct {
	TimeShift int64 `form:"time_shift" json:"time_shift" binding:"required"`
	Duration  int64 `form:"duration" json:"duration" binding:"required"`
}

func (f *GetTrafficForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}
