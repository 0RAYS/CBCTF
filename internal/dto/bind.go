package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

// Validator is an optional interface that DTOs can implement
// to add custom validation logic beyond ShouldBind.
type Validator interface {
	Validate(ctx *gin.Context) model.RetVal
}

// Bind binds the request to the given form struct and runs validation.
// If the form implements Validator, its Validate method is called after binding.
func Bind[T any](ctx *gin.Context, form *T) model.RetVal {
	if err := ctx.ShouldBind(form); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	if v, ok := any(form).(Validator); ok {
		return v.Validate(ctx)
	}
	return model.SuccessRetVal()
}
