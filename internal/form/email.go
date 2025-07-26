package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"github.com/gin-gonic/gin"
)

// VerifyEmail 邮箱验证表单
type VerifyEmail struct {
	ID    string `form:"id" binding:"required"`
	Token string `form:"token" binding:"required"`
}

func (f *VerifyEmail) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}
