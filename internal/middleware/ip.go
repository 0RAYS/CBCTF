package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
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

	ctx.Next()

	statusCode := ctx.Writer.Status()

	log := model.IP{
		IP:        ip,
		Time:      accessTime,
		Method:    method,
		URL:       url,
		UserAgent: userAgent,
		Status:    statusCode,
		Referer:   referer,
	}
	tx := db.DB.WithContext(ctx).Begin()
	db.RecordIP(tx, log)
	tx.Commit()
}

// 频率限制的配置
const (
	RateLimitWindow      = 1 * time.Minute
	RateLimitMaxRequests = 100
)

// RateLimit 实现频率限制
func RateLimit() gin.HandlerFunc {
	var requestCounts = make(map[string][]time.Time)
	var mu sync.Mutex

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "::1" || ip == "127.0.0.1" {
			c.Next()
			return
		}
		now := time.Now()

		mu.Lock()
		defer mu.Unlock()

		times := requestCounts[ip]
		var validTimes []time.Time
		for _, t := range times {
			if now.Sub(t) <= RateLimitWindow {
				validTimes = append(validTimes, t)
			}
		}
		requestCounts[ip] = validTimes
		if len(requestCounts[ip]) >= RateLimitMaxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{"msg": "TooManyRequests", "data": nil})
			c.Abort()
			return
		}
		requestCounts[ip] = append(requestCounts[ip], now)
		c.Next()
	}
}
