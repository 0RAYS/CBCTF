package middleware

import (
	"CBCTF/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
func Cors(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", config.Env.Frontend) // 可将将 * 替换为指定的域名
	ctx.Header("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, X-M, Connection, Upgrade, Sec-Websocket-Extensions, Sec-Websocket-Key, Sec-Websocket-Version")
	ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Authorization")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}
