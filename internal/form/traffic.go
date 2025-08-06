package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"

	"github.com/gin-gonic/gin"
)

type GetTrafficForm struct {
	TimeShift int64 `form:"time_shift" json:"time_shift"`
	Duration  int64 `form:"duration" json:"duration"`
}

func (f *GetTrafficForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	if f.Duration < 1 {
		f.Duration = 60
	}
	return true, i18n.Success
}
