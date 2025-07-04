package form

import (
	"CBCTF/internel/i18n"
	"github.com/gin-gonic/gin"
)

// DeleteFileForm for delete files
type DeleteFileForm struct {
	FileIDL []string `form:"file_id" json:"file_id" binding:"required"`
}

func (f *DeleteFileForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}
