package form

import (
	"CBCTF/internel/i18n"
	"github.com/gin-gonic/gin"
	"slices"
	"strings"
)

var allowedModel = []string{"user", "team", "contest", "challenge"}

type SearchForm struct {
	Limit  int    `form:"limit" json:"limit"`
	Offset int    `form:"offset" json:"offset"`
	Model  string `form:"model" json:"model" binding:"required"`
}

func (f *SearchForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	f.Model = strings.TrimSpace(strings.ToLower(f.Model))
	if !slices.Contains(allowedModel, f.Model) {
		return false, i18n.BadRequest
	}
	if f.Limit > 100 {
		f.Limit = 15
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		f.Limit = 10
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		f.Offset = 0
	}
	return true, i18n.Success
}

// GetModelsForm for get models list
type GetModelsForm struct {
	Offset int `form:"offset" json:"offset"`
	Limit  int `form:"limit" json:"limit"`
}

func (f *GetModelsForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	if f.Limit > 100 {
		f.Limit = 15
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		f.Limit = 10
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		f.Offset = 0
	}
	return true, i18n.Success
}

// ChangePasswordForm for user or admin change password
type ChangePasswordForm struct {
	OldPassword string `form:"old" json:"old" binding:"required"`
	NewPassword string `form:"new" json:"new" binding:"required"`
}

func (f *ChangePasswordForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	if f.OldPassword == f.NewPassword {
		return false, i18n.PasswordSame
	}
	return true, i18n.Success
}
