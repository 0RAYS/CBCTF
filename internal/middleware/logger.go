package middleware

import (
	"CBCTF/internal/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

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
		if raw != "" {
			path = path + "?" + raw
		}
		e := l.WithFields(logrus.Fields{
			"Latency":    n.Sub(start),
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
