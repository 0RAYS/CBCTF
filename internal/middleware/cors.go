package middleware

import (
	"CBCTF/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{config.Env.Frontend},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "X-M", "Connection", "Upgrade",
			"Sec-Websocket-Extensions", "Sec-Websocket-Key", "Sec-Websocket-Version",
		},
		ExposeHeaders: []string{
			"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control",
			"Content-Language", "Content-Type", "Authorization",
		},
		AllowCredentials: true,
	})
}
