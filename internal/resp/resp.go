package resp

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	CTXStatusCodeKey = "StatusCode"
)

func JSON(ctx *gin.Context, ret model.RetVal) {
	code, err := strconv.Atoi(i18n.Translate("und", ret.Msg))
	if err != nil {
		code = 500
	}
	ctx.Set(CTXStatusCodeKey, code)
	ctx.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   i18n.Translate(i18n.DetectLanguage(ctx), ret.Msg, ret.Attr),
		"data":  ret.Data,
		"trace": ctx.GetString("TraceID"),
	})
}

func AbortJSON(ctx *gin.Context, ret model.RetVal) {
	ctx.Abort()
	JSON(ctx, ret)
}
