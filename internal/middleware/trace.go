package middleware

import (
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
)

// Trace 设置 trace, 方便追踪日志
func Trace(ctx *gin.Context) {
	ctx.Set("TraceID", utils.RandomString())
}
