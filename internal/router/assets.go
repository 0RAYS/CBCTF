package router

import (
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/oauth"
	"net/http"

	"github.com/gin-gonic/gin"
)

var DefaultPicture = map[string][]byte{
	"github":       oauth.GithubMark,
	"github-white": oauth.GithubMarkWhite,
	"hduhelp":      oauth.HDUHelpPicture,
}

func DefaultAssets(ctx *gin.Context) {
	var form dto.GetAssetForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	file, ok := DefaultPicture[form.Filename]
	if !ok {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Model.File.NotFound})
		return
	}
	ctx.Data(http.StatusOK, "application/octet-stream", file)
}
