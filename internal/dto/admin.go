package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"

	"github.com/gin-gonic/gin"
)

// CreateAdminForm for create admin
type CreateAdminForm RegisterForm

func (f *CreateAdminForm) Validate(_ *gin.Context) model.RetVal {
	if utils.CheckPassword(f.Password) < 2 {
		return model.RetVal{Msg: i18n.Model.User.WeakPassword}
	}
	return model.SuccessRetVal()
}

// UpdateAdminForm for admin update info
type UpdateAdminForm struct {
	Name  *string `form:"name" json:"name" binding:"omitempty,min=1"`
	Email *string `form:"email" json:"email" binding:"omitempty,email"`
}
