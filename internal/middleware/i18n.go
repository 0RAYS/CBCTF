package middleware

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	CTXStatusCodeKey = "StatusCode"
)

type i18nResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *i18nResponseWriter) Write(p []byte) (n int, err error) {
	return w.body.Write(p)
}

func I18n(ctx *gin.Context) {
	w := &i18nResponseWriter{
		ResponseWriter: ctx.Writer,
		body:           bytes.NewBufferString(""),
	}
	ctx.Writer = w

	ctx.Next()

	ctx.Set(CTXStatusCodeKey, ctx.Writer.Status())
	var old struct {
		Msg  string         `json:"msg"`
		Data any            `json:"data"`
		Attr map[string]any `json:"attr"`
	}
	data := w.body.String()
	if err := json.Unmarshal([]byte(data), &old); err != nil {
		_, _ = w.ResponseWriter.Write([]byte(data))
		return
	}
	code, err := strconv.Atoi(i18n.Translate("und", old.Msg))
	if err != nil {
		log.Logger.Warningf("i18n.Translate err: %s", err)
		code = 500
	}
	ctx.Set(CTXStatusCodeKey, code)
	var res = struct {
		Code  int    `json:"code"`
		Msg   string `json:"msg"`
		Data  any    `json:"data"`
		Trace string `json:"trace"`
	}{
		Code:  code,
		Msg:   i18n.Translate(detectLanguage(ctx), old.Msg, old.Attr),
		Data:  old.Data,
		Trace: GetTraceID(ctx),
	}
	ret, err := json.Marshal(res)
	if err != nil {
		log.Logger.Errorf("Rewrite response error: %s", err)
		return
	}
	defer w.body.Reset()
	_, _ = w.ResponseWriter.Write(ret)
}

func detectLanguage(ctx *gin.Context) string {
	if lang := ctx.Query("lang"); lang != "" {
		return lang
	}
	return ctx.GetHeader("Accept-Language")
}
