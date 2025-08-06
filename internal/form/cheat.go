package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"slices"

	"github.com/gin-gonic/gin"
)

var allowedCheatType = []string{model.Suspicious, model.Cheater}

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
	if f.Limit > 100 {
		f.Limit = 15
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
