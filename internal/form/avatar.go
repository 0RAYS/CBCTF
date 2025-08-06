package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"

	"github.com/gin-gonic/gin"
)

// DeleteFileForm for delete files
type DeleteFileForm struct {
	FileIDL []string `form:"file_ids" json:"file_ids" binding:"required"`
}

func (f *DeleteFileForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}
