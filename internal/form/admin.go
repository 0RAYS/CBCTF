package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

// CreateAdminForm for create admin
type CreateAdminForm RegisterForm

func (f *CreateAdminForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	if utils.CheckPassword(f.Password) < 2 {
		return false, i18n.WeakPassword
	}
	f.Name = strings.TrimSpace(f.Name)
	if f.Name == "" {
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}

// UpdateAdminForm for admin update info
type UpdateAdminForm struct {
	Name  *string `form:"name" json:"name"`
	Email *string `form:"email" json:"email" binding:"omitempty,email"`
}

func (f *UpdateAdminForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	if f.Name != nil {
		*f.Name = strings.TrimSpace(*f.Name)
		if *f.Name == "" {
			return false, i18n.BadRequest
		}
	}
	return true, i18n.Success
}
