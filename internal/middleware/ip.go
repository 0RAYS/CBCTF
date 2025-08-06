package middleware

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/redis"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit(name string, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		target := ctx.ClientIP()
		if slices.Contains(config.Env.Gin.RateLimit.Whitelist, target) {
			ctx.Next()
			return
		}
		if userID := GetSelfID(ctx); userID != 0 {
			target = fmt.Sprintf("%d", userID)
		}
		count, err := redis.RateLimit(name, target, window)
		if err != nil {
			log.Logger.Warningf("Failed to rate limit: %s", err)
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.RedisError, "data": nil})
			return
		}
		if int(count) > maxRequests {
			go prometheus.UpdateRateLimitMetrics(name, target)
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.TooManyRequests, "data": nil})
			return
		}
		ctx.Next()
	}
}
