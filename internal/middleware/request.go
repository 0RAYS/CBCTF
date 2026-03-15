package middleware

import (
	"CBCTF/internal/config"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"database/sql"
	"slices"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

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

// AccessLog 记录访问日志
func AccessLog(ctx *gin.Context) {
	ip := ctx.ClientIP()
	accessTime := time.Now()
	method := ctx.Request.Method
	url := ctx.Request.URL.Path
	userAgent := ctx.Request.UserAgent()
	referer := ctx.Request.Referer()
	magic := GetMagic(ctx)
	path := ctx.FullPath()
	ctx.Next()

	statusCode := ctx.GetInt(resp.CTXStatusCodeKey)
	selfID := GetSelf(ctx).ID

	if !slices.Contains(config.Env.Gin.Log.Whitelist, path) {
		// Truncate long headers to 255 characters to fit storage constraints
		if len(userAgent) > 255 {
			userAgent = userAgent[:255]
		}
		if len(referer) > 255 {
			referer = referer[:255]
		}
		request := model.Request{
			IP:        ip,
			Time:      accessTime,
			Method:    method,
			Path:      path,
			URL:       url,
			UserAgent: userAgent,
			Status:    statusCode,
			Referer:   referer,
			Magic:     magic,
		}
		if selfID > 0 {
			request.UserID = sql.Null[uint]{V: selfID, Valid: true}
		}
		AppendRequest(request)
	}
}
