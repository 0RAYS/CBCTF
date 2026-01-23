package middleware

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/redis"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit(name string, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		client := ctx.ClientIP()
		if slices.Contains(config.Env.Gin.RateLimit.Whitelist, client) {
			ctx.Next()
			return
		}
		if userID := GetSelfID(ctx); userID != 0 {
			client = strconv.Itoa(int(userID))
		}
		count, ret := redis.RateLimit(name, client, window)
		if !ret.OK {
			ctx.AbortWithStatusJSON(http.StatusOK, ret)
			return
		}
		if int(count) > maxRequests {
			prometheus.UpdateRateLimitMetrics(name, client)
			ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Request.TooManyRequests})
			return
		}
		ctx.Next()
	}
}
