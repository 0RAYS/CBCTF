package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"slices"

	"github.com/gin-gonic/gin"
)

var allowedFileType = []string{model.WriteUPFileType, model.AvatarFileType, model.ChallengeFileType, model.TrafficFileType}

type GetFilesForm struct {
	Offset int    `form:"offset" json:"offset"`
	Limit  int    `form:"limit" json:"limit"`
	Type   string `form:"type" json:"type"`
}

func (f *GetFilesForm) Bind(ctx *gin.Context) (bool, string) {
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
	if !slices.Contains(allowedFileType, f.Type) {
		f.Type = ""
	}
	return true, i18n.Success
}
