package middleware

import (
	"CBCTF/internal/utils"

	"github.com/gin-gonic/gin"
)

// SetTrace 设置 trace, 方便追踪日志
func SetTrace(ctx *gin.Context) {
	ctx.Set("TraceID", utils.UUID())
}

// GetTraceID 从 gin.Context 中获取 trace, 该值由 middleware.Trace 设置
func GetTraceID(ctx *gin.Context) string {
	return ctx.GetString("TraceID")
}
