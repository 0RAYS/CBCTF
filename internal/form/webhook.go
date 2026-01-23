package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

var allowMethods = []string{http.MethodPost, http.MethodGet}

type CreateWebhookForm struct {
	Name    string           `form:"name" json:"name" binding:"required"`
	URL     string           `form:"url" json:"url" binding:"required"`
	Method  string           `form:"method" json:"method" binding:"required"`
	Headers model.StringMap  `form:"headers" json:"headers"`
	Timeout int64            `form:"timeout" json:"timeout"`
	Retry   int              `form:"retry" json:"retry"`
	Events  model.StringList `gorm:"type:json" json:"events"`
}

func (f *CreateWebhookForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	f.Method = strings.ToUpper(f.Method)
	if !slices.Contains(allowMethods, f.Method) {
		return model.RetVal{Msg: i18n.Model.Webhook.InvalidMethod}
	}
	return model.SuccessRetVal()
}

type UpdateWebhookForm struct {
	Name    *string           `form:"name" json:"name"`
	URL     *string           `form:"url" json:"url"`
	Method  *string           `form:"method" json:"method"`
	Headers *model.StringMap  `form:"headers" json:"headers"`
	Timeout *int64            `form:"timeout" json:"timeout"`
	Retry   *int              `form:"retry" json:"retry"`
	On      *bool             `form:"on" json:"on"`
	Events  *model.StringList `form:"events" json:"events"`
}

func (f *UpdateWebhookForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	if f.Method != nil {
		*f.Method = strings.ToUpper(*f.Method)
		if !slices.Contains(allowMethods, *f.Method) {
			return model.RetVal{Msg: i18n.Model.Webhook.InvalidMethod}
		}
	}
	return model.SuccessRetVal()
}
