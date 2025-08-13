package middleware

import (
	"CBCTF/internal/log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var TotalDuration time.Duration
var TotalRequests int
var MU sync.Mutex

func Logger(ctx *gin.Context) {
	l := log.Logger.WithField("type", "GIN")
	start := time.Now()
	path := ctx.Request.URL.Path
	raw := ctx.Request.URL.RawQuery

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
	if raw != "" {
		path = path + "?" + raw
	}
	e := l.WithFields(logrus.Fields{
		"Latency":    latency,
		"StatusCode": ctx.GetInt(CTXStatusCodeKey),
		"Method":     ctx.Request.Method,
		"ClientIP":   ctx.ClientIP(),
		"Path":       path,
		"TraceID":    GetTraceID(ctx),
	})

	if ctx.Errors != nil {
		e.Error(ctx.Errors.String())
	} else {
		e.Info()
	}
}
