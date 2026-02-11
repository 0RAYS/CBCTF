package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type CreateSmtpForm struct {
	Address string `form:"address" json:"address" binding:"required,email"`
	Host    string `form:"host" json:"host" binding:"required,hostname"`
	Port    int    `form:"port" json:"port" binding:"required,gte=0,lte=65535"`
	Pwd     string `form:"pwd" json:"pwd" binding:"required"`
}

func (f *CreateSmtpForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

type UpdateSmtpForm struct {
	Address *string `form:"address" json:"address" binding:"omitempty,email"`
	Host    *string `form:"host" json:"host" binding:"omitempty,hostname"`
	Port    *int    `form:"port" json:"port" binding:"omitempty,gte=0,lte=65535"`
	Pwd     *string `form:"pwd" json:"pwd"`
	On      *bool   `form:"on" json:"on"`
}

func (f *UpdateSmtpForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
