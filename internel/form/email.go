package form

import (
	"CBCTF/internel/i18n"
	"github.com/gin-gonic/gin"
)

// VerifyEmail 邮箱验证表单
type VerifyEmail struct {
	ID    string `form:"id" binding:"required"`
	Token string `form:"token" binding:"required"`
}

func (f *VerifyEmail) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}
