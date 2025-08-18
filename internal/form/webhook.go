package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type CreateWebhookForm struct {
	Name       string           `form:"name" json:"name" binding:"required"`
	URL        string           `form:"url" json:"url" binding:"required"`
	Method     string           `form:"method" json:"method" binding:"required"`
	Headers    model.StringMap  `form:"headers" json:"headers"`
	Timeout    int64            `form:"timeout" json:"timeout"`
	RetryCount int64            `form:"retry_count" json:"retry_count"`
	RetryDelay int64            `form:"retry_delay" json:"retry_delay"`
	Events     model.StringList `gorm:"type:json" json:"events"`
}

func (f *CreateWebhookForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}

type UpdateWebhookForm struct {
	Name       *string           `form:"name" json:"name"`
	URL        *string           `form:"url" json:"url"`
	Method     *string           `form:"method" json:"method"`
	Headers    *model.StringMap  `form:"headers" json:"headers"`
	Timeout    *int64            `form:"timeout" json:"timeout"`
	RetryCount *int64            `form:"retry_count" json:"retry_count"`
	RetryDelay *int64            `form:"retry_delay" json:"retry_delay"`
	On         *bool             `form:"on" json:"on"`
	Events     *model.StringList `form:"events" json:"events"`
}

func (f *UpdateWebhookForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}
