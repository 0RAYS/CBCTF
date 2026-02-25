package middleware

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/redis"
	"CBCTF/internal/resp"
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
		if userID := GetSelf(ctx).ID; userID != 0 {
			client = strconv.Itoa(int(userID))
		}
		count, ret := redis.RateLimit(name, client, window)
		if !ret.OK {
			resp.AbortJSON(ctx, ret)
			return
		}
		if int(count) > maxRequests {
			prometheus.RecordRateLimitHit(name)
			resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.TooManyRequests})
			return
		}
		ctx.Next()
	}
}
