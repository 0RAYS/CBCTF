package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

// CreateTeamForm for create team
type CreateTeamForm struct {
	Name    string `form:"name" json:"name" binding:"required"`
	Desc    string `form:"desc" json:"desc"`
	Captcha string `form:"captcha" json:"captcha"`
}

func (f *CreateTeamForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	f.Name = strings.TrimSpace(f.Name)
	if f.Name == "" {
		return false, i18n.BadRequest
	}
	f.Captcha = strings.TrimSpace(f.Captcha)
	return true, i18n.Success
}

// UpdateTeamForm for user update team info
type UpdateTeamForm struct {
	Name      *string `form:"name" json:"name"`
	Desc      *string `form:"desc" json:"desc"`
	CaptainID *uint   `form:"captain_id" json:"captain_id"`
}

func (f *UpdateTeamForm) Bind(ctx *gin.Context) (bool, string) {
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

// JoinTeamForm for user join team
type JoinTeamForm struct {
	Name    string `form:"name" json:"name" binding:"required"`
	Captcha string `form:"captcha" json:"captcha" binding:"required"`
}

func (f *JoinTeamForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	f.Name = strings.TrimSpace(f.Name)
	if f.Name == "" {
		return false, i18n.BadRequest
	}
	f.Captcha = strings.TrimSpace(f.Captcha)
	return true, i18n.Success
}

// KickMemberForm for admin or captain kick member
type KickMemberForm struct {
	UserID uint `form:"user_id" json:"user_id" binding:"required"`
}

func (f *KickMemberForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}

// AdminUpdateTeamForm for admin update team info
type AdminUpdateTeamForm struct {
	Name      *string `form:"name" json:"name"`
	Desc      *string `form:"desc" json:"desc"`
	Hidden    *bool   `form:"hidden" json:"hidden"`
	Banned    *bool   `form:"banned" json:"banned"`
	Captcha   *string `form:"captcha" json:"captcha"`
	CaptainID *uint   `form:"captain_id" json:"captain_id"`
}

func (f *AdminUpdateTeamForm) Bind(ctx *gin.Context) (bool, string) {
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
	if f.Captcha != nil {
		f.Captcha = utils.Ptr(strings.TrimSpace(*f.Captcha))
	}
	return true, i18n.Success
}
