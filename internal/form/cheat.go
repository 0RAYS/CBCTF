package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

var allowedCheatType = []string{model.Suspicious, model.Cheater, model.Pass}

type GetCheatsForm struct {
	Offset int    `form:"offset" json:"offset"`
	Limit  int    `form:"limit" json:"limit"`
	Type   string `form:"type" json:"type"`
}

func (f *GetCheatsForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	if f.Limit > 100 || f.Limit < 0 {
		f.Limit = 15
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	if !slices.Contains(allowedCheatType, f.Type) {
		f.Type = ""
	}
	return true, i18n.Success
}

type UpdateCheatForm struct {
	Reason  *string `form:"reason" json:"reason"`
	Type    *string `form:"type" json:"type"`
	Checked *bool   `form:"checked" json:"checked"`
	Comment *string `form:"comment" json:"comment"`
}

func (f *UpdateCheatForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBindJSON(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	if f.Type != nil {
		f.Type = utils.Ptr(strings.Trim(*f.Type, " "))
		if !slices.Contains(allowedCheatType, *f.Type) {
			return false, i18n.BadRequest
		}
	}
	return true, i18n.Success
}
