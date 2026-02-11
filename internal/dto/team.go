package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

// CreateTeamForm for create team
type CreateTeamForm struct {
	Name        string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description"`
	Captcha     string `form:"captcha" json:"captcha"`
}

func (f *CreateTeamForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

// UpdateTeamForm for user update team info
type UpdateTeamForm struct {
	Name        *string `form:"name" json:"name" binding:"omitempty,min=1"`
	Description *string `form:"description" json:"description"`
	CaptainID   *uint   `form:"captain_id" json:"captain_id"`
}

func (f *UpdateTeamForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

// JoinTeamForm for user join team
type JoinTeamForm struct {
	Name    string `form:"name" json:"name" binding:"required"`
	Captcha string `form:"captcha" json:"captcha" binding:"required"`
}

func (f *JoinTeamForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

// KickMemberForm for admin or captain kick member
type KickMemberForm struct {
	UserID uint `form:"user_id" json:"user_id" binding:"required"`
}

func (f *KickMemberForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

// AdminUpdateTeamForm for admin update team info
type AdminUpdateTeamForm struct {
	Name        *string `form:"name" json:"name" binding:"omitempty,min=1"`
	Description *string `form:"description" json:"description"`
	Hidden      *bool   `form:"hidden" json:"hidden"`
	Banned      *bool   `form:"banned" json:"banned"`
	Captcha     *string `form:"captcha" json:"captcha" binding:"omitempty,min=1"`
	CaptainID   *uint   `form:"captain_id" json:"captain_id"`
}

func (f *AdminUpdateTeamForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
