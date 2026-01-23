package router

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/oauth"
	"net/http"

	"github.com/gin-gonic/gin"
)

var DefaultAvatar = map[string][]byte{
	"github":       oauth.GithubMark,
	"github-white": oauth.GithubMarkWhite,
	"hduhelp":      oauth.HDUHelpAvatar,
}

func DefaultAssets(ctx *gin.Context) {
	var form f.GetAssetForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	file, ok := DefaultAvatar[form.Filename]
	if !ok {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Model.NotFound, Attr: map[string]any{"Model": model.File{}.GetModelName()}})
		return
	}
	ctx.Data(http.StatusOK, "application/octet-stream", file)
}
