package middleware

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				panicErr := fmt.Errorf("panic recovered: %v", err)
				_ = ctx.Error(panicErr)

				log.Logger.WithFields(logrus.Fields{
					"Type":     log.GinLogType,
					"Method":   ctx.Request.Method,
					"Path":     ctx.Request.URL.RequestURI(),
					"ClientIP": ctx.ClientIP(),
					"TraceID":  GetTraceID(ctx),
					"Stack":    string(debug.Stack()),
				}).Error(panicErr)

				resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": fmt.Sprint(err)}})
			}
		}()

		ctx.Next()
	}
}
