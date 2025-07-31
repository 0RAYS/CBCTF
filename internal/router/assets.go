package router

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/oauth"
	"github.com/gin-gonic/gin"
	"net/http"
)

var DefaultAvatar = map[string][]byte{
	"github":       oauth.GithubMark,
	"github-white": oauth.GithubMarkWhite,
}

func DefaultAssets(ctx *gin.Context) {
	var form f.GetAssetForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	file, ok := DefaultAvatar[form.Filename]
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.FileNotFound, "data": nil})
		return
	}
	ctx.Data(http.StatusOK, "application/octet-stream", file)
}
