package middleware

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	db "CBCTF/internel/repo"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
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

	if path != "/metrics" {
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

var requestCounts = make(map[string][]time.Time)
var mu sync.Mutex

// RateLimit 实现频率限制
func RateLimit(ctx *gin.Context) {
	ip := ctx.ClientIP()
	if ip == "::1" || ip == "127.0.0.1" {
		ctx.Next()
		return
	}
	now := time.Now()

	mu.Lock()
	defer mu.Unlock()

	times := requestCounts[ip]
	var validTimes []time.Time
	for _, t := range times {
		if now.Sub(t) <= time.Duration(config.Env.Gin.RateLimit.Window)*time.Second {
			validTimes = append(validTimes, t)
		}
	}
	requestCounts[ip] = validTimes
	if len(requestCounts[ip]) >= config.Env.Gin.RateLimit.MaxRequests {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.TooManyRequests, "data": nil})
		ctx.Abort()
		return
	}
	requestCounts[ip] = append(requestCounts[ip], now)
	ctx.Next()
}
