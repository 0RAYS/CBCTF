package middleware

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/resp"
	"slices"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var TotalDuration time.Duration
var TotalRequests int
var MU sync.Mutex

func Logger(ctx *gin.Context) {
	start := time.Now()

	// Process request
	ctx.Next()
	// Stop timer
	n := time.Now()
	latency := n.Sub(start)
	if ctx.Request.Method != "OPTIONS" {
		MU.Lock()
		TotalDuration += latency
		TotalRequests++
		MU.Unlock()
	}
	path := ctx.Request.URL.Path
	raw := ctx.Request.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}
	if slices.Contains(config.Env.Gin.Log.Whitelist, ctx.FullPath()) {
		return
	}
	statusCode := ctx.GetInt(resp.CTXStatusCodeKey)
	if statusCode == 0 {
		statusCode = ctx.Writer.Status()
	}
	e := log.Logger.WithFields(logrus.Fields{
		"Type":       log.GinLogType,
		"Latency":    latency,
		"StatusCode": statusCode,
		"Method":     ctx.Request.Method,
		"ClientIP":   ctx.ClientIP(),
		"Path":       path,
		"TraceID":    GetTraceID(ctx),
	})

	if ctx.Errors != nil {
		e.Error(ctx.Errors.String())
	} else if statusCode >= 500 {
		e.Error()
	} else if statusCode >= 400 {
		e.Warning()
	} else {
		e.Info()
	}
}
