package middleware

import (
	"CBCTF/internal/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var TotalDuration time.Duration
var TotalRequests int
var MU sync.Mutex

func Logger() func(ctx *gin.Context) {
	l := log.Logger.WithField("type", "GIN")
	return func(r *gin.Context) {
		// Start timer
		start := time.Now()
		path := r.Request.URL.Path
		raw := r.Request.URL.RawQuery

		// Process request
		r.Next()
		// Stop timer
		n := time.Now()
		latency := n.Sub(start)
		MU.Lock()
		TotalDuration += latency
		TotalRequests++
		MU.Unlock()
		if raw != "" {
			path = path + "?" + raw
		}
		e := l.WithFields(logrus.Fields{
			"Latency":    latency,
			"StatusCode": r.Writer.Status(),
			"Method":     r.Request.Method,
			"ClientIP":   r.ClientIP(),
			"Path":       path,
			"TraceID":    r.Value("TraceID"),
		})

		if r.Errors != nil {
			e.Error(r.Errors.String())
		} else {
			e.Info()
		}
	}
}
