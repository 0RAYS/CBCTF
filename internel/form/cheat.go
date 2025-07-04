package form

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
	"slices"
)

var allowedCheatType = []string{model.Suspicious, model.Cheater}

type GetCheatsForm struct {
	Offset int    `form:"offset" json:"offset"`
	Limit  int    `form:"limit" json:"limit"`
	Type   string `form:"type" json:"type"`
}

func (f *GetCheatsForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		f.Limit = 10
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		f.Offset = 0
	}
	if !slices.Contains(allowedCheatType, f.Type) {
		f.Type = ""
	}
	return true, i18n.Success
}
