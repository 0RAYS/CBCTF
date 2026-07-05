package middleware

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"database/sql"
	"net/http"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var TotalDuration atomic.Int64
var TotalRequests atomic.Int64

var (
	RequestsPool  = make([]model.Request, 0)
	RequestsMutex sync.Mutex
)

func AppendRequest(request model.Request) {
	RequestsMutex.Lock()
	RequestsPool = append(RequestsPool, request)
	RequestsMutex.Unlock()
}

func DrainRequestsPool() []model.Request {
	RequestsMutex.Lock()

	if len(RequestsPool) == 0 {
		RequestsMutex.Unlock()
		return nil
	}

	requests := RequestsPool
	RequestsPool = make([]model.Request, 0)
	RequestsMutex.Unlock()
	return requests
}

func Logger(ctx *gin.Context) {
	if slices.Contains(config.Env.Gin.Log.Whitelist, ctx.FullPath()) {
		ctx.Next()
		return
	}
	start := time.Now()

	// Process request
	ctx.Next()
	// Stop timer
	latency := time.Now().Sub(start)
	if ctx.Request.Method != "OPTIONS" {
		TotalDuration.Add(latency.Milliseconds())
		TotalRequests.Add(1)
	}
	path := ctx.Request.URL.Path
	raw := ctx.Request.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}
	statusCode := ctx.GetInt(resp.CTXStatusCodeKey)
	if statusCode == 0 {
		statusCode = ctx.Writer.Status()
	}
	selfID := GetSelf(ctx).ID
	userAgent := ctx.Request.UserAgent()
	if len(userAgent) > 255 {
		userAgent = userAgent[:255]
	}
	referer := ctx.Request.Referer()
	if len(referer) > 255 {
		referer = referer[:255]
	}

	AppendRequest(model.Request{
		IP:        ctx.ClientIP(),
		Time:      start,
		Latency:   latency,
		Method:    ctx.Request.Method,
		Path:      ctx.FullPath(),
		URL:       ctx.Request.URL.String(),
		UserAgent: userAgent,
		Status:    statusCode,
		Referer:   referer,
		UserID:    sql.Null[uint]{V: selfID, Valid: selfID > 0},
	})

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
	} else if statusCode >= http.StatusInternalServerError {
		e.Error()
	} else if statusCode >= http.StatusBadRequest {
		e.Warning()
	} else {
		e.Info()
	}
}
