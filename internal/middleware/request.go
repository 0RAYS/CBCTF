package middleware

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/resp"
	"database/sql"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
)

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
		request := db.CreateRequestOptions{
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
		db.InitRequestRepo(db.DB).Create(request)
	}
}
