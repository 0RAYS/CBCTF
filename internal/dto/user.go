package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"

	"github.com/gin-gonic/gin"
)

// LoginForm for user or admin login
type LoginForm struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// RegisterForm for user register
type RegisterForm struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
}

func (f *RegisterForm) Validate(ctx *gin.Context) model.RetVal {
	if utils.CheckPassword(f.Password) < 2 {
		return model.RetVal{Msg: i18n.Model.User.WeakPassword}
	}
	return model.SuccessRetVal()
}

// CreateUserForm for create user
type CreateUserForm struct {
	Name        string `form:"name" json:"name" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	Email       string `form:"email" json:"email" binding:"required,email"`
	Description string `form:"description" json:"description"`
	Hidden      bool   `form:"hidden" json:"hidden"`
	Verified    bool   `form:"verified" json:"verified"`
	Banned      bool   `form:"banned" json:"banned"`
}

func (f *CreateUserForm) Validate(ctx *gin.Context) model.RetVal {
	if utils.CheckPassword(f.Password) < 2 {
		return model.RetVal{Msg: i18n.Model.User.WeakPassword}
	}
	return model.SuccessRetVal()
}

// UpdateSelfForm for user update info
type UpdateSelfForm struct {
	Name        *string `form:"name" json:"name" binding:"omitempty,min=1"`
	Email       *string `form:"email" json:"email" binding:"omitempty,email"`
	Description *string `form:"description" json:"description"`
}

// UpdateUserForm for admin update user info
type UpdateUserForm struct {
	Name        *string `form:"name" json:"name" binding:"omitempty,min=1"`
	Email       *string `form:"email" json:"email" binding:"omitempty,email"`
	Description *string `form:"description" json:"description"`
	Password    *string `form:"password" json:"password" binding:"omitempty,min=1"`
	Hidden      *bool   `form:"hidden" json:"hidden"`
	Banned      *bool   `form:"banned" json:"banned"`
	Verified    *bool   `form:"verified" json:"verified"`
}

func (f *UpdateUserForm) Validate(ctx *gin.Context) model.RetVal {
	if f.Password != nil {
		if utils.CheckPassword(*f.Password) < 2 {
			return model.RetVal{Msg: i18n.Model.User.WeakPassword}
		}
	}
	return model.SuccessRetVal()
}

// DeleteSelfForm for user delete self
type DeleteSelfForm struct {
	Password string `form:"password" json:"password" binding:"required"`
}
