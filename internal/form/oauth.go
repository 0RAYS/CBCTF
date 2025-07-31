package form

import (
	"CBCTF/internal/i18n"
	"github.com/gin-gonic/gin"
)

type OauthCallbackForm struct {
	Code  string `form:"code" json:"code" binding:"required"`
	State string `form:"state" json:"state" binding:"required"`
}

func (f *OauthCallbackForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}
