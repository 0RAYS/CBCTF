package middleware

import (
	"github.com/gin-gonic/gin"
)

// SetMagic 保存设备 ID 至上下文
func SetMagic(ctx *gin.Context) {
	ctx.Set("Magic", ctx.GetHeader("X-M"))
	ctx.Next()
}

// GetMagic 获取设备 ID
func GetMagic(ctx *gin.Context) string {
	return ctx.GetString("Magic")
}
