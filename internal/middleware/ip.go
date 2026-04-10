package middleware

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/ratelimit"
	"CBCTF/internal/resp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit(name string, maxRequests int, window time.Duration) gin.HandlerFunc {
	rule := ratelimit.Rule{
		Name:   name,
		Limit:  int64(maxRequests),
		Window: window,
	}
	return func(ctx *gin.Context) {
		subject := "ip:" + ctx.ClientIP()
		if userID := GetSelf(ctx).ID; userID != 0 {
			subject = "user:" + strconv.Itoa(int(userID))
		}
		decision, err := ratelimit.RateLimiter.Allow(ctx.Request.Context(), rule, ratelimit.Subject{
			Key:      subject,
			ClientIP: ctx.ClientIP(),
		})
		if err != nil {
			resp.AbortJSON(ctx, model.RetVal{
				Msg:  i18n.Redis.SetError,
				Attr: map[string]any{"Key": name, "Error": err.Error()},
			})
			return
		}
		if !decision.Allowed {
			retryAfterSeconds := (decision.RetryAfter.Milliseconds() + 999) / 1000
			if retryAfterSeconds <= 0 {
				retryAfterSeconds = 1
			}
			prometheus.RecordRateLimitHit(name)
			resp.AbortJSON(ctx, model.RetVal{
				Msg:  i18n.Response.TooManyRequests,
				Attr: map[string]any{"Seconds": retryAfterSeconds},
			})
			return
		}
		ctx.Next()
	}
}
