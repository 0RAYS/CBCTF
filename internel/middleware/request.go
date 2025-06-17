package middleware

import (
	"CBCTF/internel/config"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"github.com/gin-gonic/gin"
	"time"
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

	if !utils.In(path, config.Env.Gin.Log.Whitelist) {
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
		tx := db.DB.WithContext(ctx).Begin()
		_, ok, _ := db.InitRequestRepo(tx).Create(request)
		if !ok {
			tx.Rollback()
			return
		}
		tx.Commit()
	}
}
