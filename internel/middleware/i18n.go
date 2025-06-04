package middleware

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strings"
)

type i18nResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *i18nResponseWriter) Write(p []byte) (n int, err error) {
	return w.body.Write(p)
}

type Data struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Data  any    `json:"data"`
	Trace string `json:"trace"`
}

// I18n 重写响应
func I18n(ctx *gin.Context) {
	w := &i18nResponseWriter{
		ResponseWriter: ctx.Writer,
		body:           bytes.NewBufferString(""),
	}
	ctx.Writer = w

	ctx.Next()

	var res Data
	old := w.body.String()

	err := json.Unmarshal([]byte(old), &res)
	if err != nil {
		_, _ = w.ResponseWriter.Write([]byte(old))
		return
	}
	language := ctx.GetHeader("Accept-Language")
	if strings.HasPrefix(language, "en-US") {
		language = "en-US"
	} else if strings.HasPrefix(language, "origin") {
		language = "origin"
	} else {
		language = "zh-CN"

	}
	res.Msg, res.Code = i18n.I18N(res.Msg, language)
	res.Trace = GetTraceID(ctx)
	ret, err := json.Marshal(res)
	if err != nil {
		log.Logger.Errorf("Rewrite response error: %v", err)
		return
	}
	defer w.body.Reset()
	_, _ = w.ResponseWriter.Write(ret)
}
