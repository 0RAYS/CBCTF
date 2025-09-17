package middleware

import (
	"CBCTF/internal/log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger(ctx *gin.Context) {
	l := log.Logger.WithField("Type", log.GinLogType)
	path := ctx.Request.URL.Path
	raw := ctx.Request.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}
	e := l.WithFields(logrus.Fields{
		"Latency":    time.Millisecond,
		"StatusCode": 200,
		"Method":     ctx.Request.Method,
		"ClientIP":   ctx.ClientIP(),
		"Path":       path,
		"TraceID":    GetTraceID(ctx),
	})
	e.Info()
	ctx.Next()
}
