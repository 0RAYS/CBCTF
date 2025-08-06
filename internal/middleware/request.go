package middleware

import (
	"CBCTF/internal/config"
	db "CBCTF/internal/repo"
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

	statusCode := ctx.Writer.Status()

	if !slices.Contains(config.Env.Gin.Log.Whitelist, path) {
		request := db.CreateRequestOptions{
			IP:        ip,
			Time:      accessTime,
			Method:    method,
			Path:      path,
			URL:       url,
			UserAgent: userAgent[:len(userAgent)%255],
			Status:    statusCode,
			Referer:   referer[:len(referer)%255],
			Magic:     magic,
		}
		tx := db.DB.WithContext(ctx).Begin()
		_, ok, _ := db.InitRequestRepo(tx).Create(request)
		if !ok {
			tx.Rollback()
			return
		}
		tx.Commit()
	}
}
