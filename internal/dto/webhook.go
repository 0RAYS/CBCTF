package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type CreateWebhookForm struct {
	Name    string           `form:"name" json:"name" binding:"required"`
	URL     string           `form:"url" json:"url" binding:"required,url"`
	Method  string           `form:"method" json:"method" binding:"required,oneof=POST GET"`
	Headers model.StringMap  `form:"headers" json:"headers"`
	Timeout int64            `form:"timeout" json:"timeout" binding:"gte=0"`
	Retry   int              `form:"retry" json:"retry" binding:"gte=0"`
	Events  model.StringList `form:"events" json:"events"`
}

func (f *CreateWebhookForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

type UpdateWebhookForm struct {
	Name    *string           `form:"name" json:"name" binding:"omitempty,min=1"`
	URL     *string           `form:"url" json:"url" binding:"omitempty,url"`
	Method  *string           `form:"method" json:"method" binding:"omitempty,oneof=POST GET"`
	Headers *model.StringMap  `form:"headers" json:"headers"`
	Timeout *int64            `form:"timeout" json:"timeout" binding:"omitempty,gte=0"`
	Retry   *int              `form:"retry" json:"retry" binding:"omitempty,gte=0"`
	On      *bool             `form:"on" json:"on"`
	Events  *model.StringList `form:"events" json:"events"`
}

func (f *UpdateWebhookForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
