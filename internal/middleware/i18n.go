package middleware

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/log"
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

func I18n() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		w := &i18nResponseWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBufferString(""),
		}
		ctx.Writer = w

		ctx.Next()

		var result Data
		old := w.body.String()

		err := json.Unmarshal([]byte(old), &result)
		if err != nil {
			_, _ = w.ResponseWriter.Write([]byte(old))
			return
		}
		language := ctx.GetHeader("Accept-Language")
		if strings.HasPrefix(language, "zh-CN") {
			language = "zh-CN"
		} else {
			language = "en-US"
		}
		result.Msg, result.Code = constants.I18N(result.Msg, language)
		result.Trace = GetTraceID(ctx)
		ret, err := json.Marshal(result)
		if err != nil {
			log.Logger.Errorf("Rewrite response error: %v", err)
			return
		}
		defer w.body.Reset()
		_, _ = w.ResponseWriter.Write(ret)
	}
}
