package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"

	"github.com/gin-gonic/gin"
)

type CreateSmtpForm struct {
	Address string `form:"address" json:"address" binding:"required"`
	Host    string `form:"host" json:"host" binding:"required"`
	Port    int    `form:"port" json:"port" binding:"required"`
	Pwd     string `form:"pwd" json:"pwd" binding:"required"`
}

func (f *CreateSmtpForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	if f.Port < 0 || f.Port > 65535 {
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}

type UpdateSmtpForm struct {
	Address *string `form:"address" json:"address"`
	Host    *string `form:"host" json:"host"`
	Port    *int    `form:"port" json:"port"`
	Pwd     *string `form:"pwd" json:"pwd"`
	On      *bool   `form:"on" json:"on"`
}

func (f *UpdateSmtpForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	if f.Port != nil {
		if *f.Port < 0 || *f.Port > 65535 {
			return false, i18n.BadRequest
		}
	}
	return true, i18n.Success
}
