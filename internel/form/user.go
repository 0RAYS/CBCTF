package form

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

// LoginForm for user or admin login
type LoginForm struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (f *LoginForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	f.Name = strings.TrimSpace(f.Name)
	if f.Name == "" {
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}

// RegisterForm for user register
type RegisterForm struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
}

func (f *RegisterForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	f.Name = strings.TrimSpace(f.Name)
	if f.Name == "" {
		return false, i18n.BadRequest
	}
	if utils.CheckPassword(f.Password) < 2 {
		return false, i18n.WeakPassword
	}
	return true, i18n.Success
}

// CreateUserForm for create user
type CreateUserForm struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Country  string `form:"country" json:"country"`
	Desc     string `form:"desc" json:"desc"`
	Hidden   bool   `form:"hidden" json:"hidden"`
	Verified bool   `form:"verified" json:"verified"`
	Banned   bool   `form:"banned" json:"banned"`
}

func (f *CreateUserForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	f.Name = strings.TrimSpace(f.Name)
	if f.Name == "" {
		return false, i18n.BadRequest
	}
	if utils.CheckPassword(f.Password) < 2 {
		return false, i18n.WeakPassword
	}
	return true, i18n.Success
}

// UpdateSelfForm for user update info
type UpdateSelfForm struct {
	Name    *string `form:"name" json:"name"`
	Email   *string `form:"email" json:"email" binding:"omitempty,email"`
	Desc    *string `form:"desc" json:"desc"`
	Country *string `form:"country" json:"country"`
}

func (f *UpdateSelfForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	if f.Name != nil {
		f.Name = utils.Ptr(strings.TrimSpace(*f.Name))
		if *f.Name == "" {
			return false, i18n.BadRequest
		}
	}
	return true, i18n.Success
}

// UpdateUserForm for admin update user info
type UpdateUserForm struct {
	Name     *string `form:"name" json:"name"`
	Email    *string `form:"name" json:"email" binding:"omitempty,email"`
	Desc     *string `form:"desc" json:"desc"`
	Country  *string `form:"country" json:"country"`
	Password *string `form:"password" json:"password"`
	Hidden   *bool   `form:"hidden" json:"hidden"`
	Banned   *bool   `form:"banned" json:"banned"`
	Verified *bool   `form:"verified" json:"verified"`
}

func (f *UpdateUserForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	if f.Name != nil {
		f.Name = utils.Ptr(strings.TrimSpace(*f.Name))
		if *f.Name == "" {
			return false, i18n.BadRequest
		}
	}
	if f.Password != nil {
		if utils.CheckPassword(*f.Password) < 2 {
			return false, i18n.WeakPassword
		}
	}
	return true, i18n.Success
}

// DeleteSelfForm for user delete self
type DeleteSelfForm struct {
	Password string `form:"password" json:"password" binding:"required"`
}

func (f *DeleteSelfForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}
