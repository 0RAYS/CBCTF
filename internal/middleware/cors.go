package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 临时使用, 后续一定要去除
func Cors(ctx *gin.Context) {
	method := ctx.Request.Method
	if ctx.Request.Header.Get("Origin") != "" {
		ctx.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, PATCH, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Authorization")
		ctx.Header("Access-Control-Allow-Credentials", "true")
	}
	if method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
	}
	ctx.Next()
}
