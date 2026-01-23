package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"

	"github.com/gin-gonic/gin"
)

// CreateAdminForm for create admin
type CreateAdminForm RegisterForm

func (f *CreateAdminForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	if utils.CheckPassword(f.Password) < 2 {
		return model.RetVal{Msg: i18n.Model.User.WeakPassword}
	}
	return model.SuccessRetVal()
}

// UpdateAdminForm for admin update info
type UpdateAdminForm struct {
	Name  *string `form:"name" json:"name"`
	Email *string `form:"email" json:"email" binding:"omitempty,email"`
}

func (f *UpdateAdminForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
