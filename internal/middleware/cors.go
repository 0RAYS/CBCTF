package middleware

import (
	"CBCTF/internal/config"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	conf := cors.Config{
		AllowOrigins: config.Env.Gin.CORS,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "X-M", "Connection", "Upgrade",
		},
		ExposeHeaders: []string{
			"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control",
			"Content-Language", "Content-Type", "Authorization", "File",
		},
		AllowCredentials: true,
	}
	if strings.ToLower(config.Env.Gin.Mode) != gin.ReleaseMode {
		conf.AllowOrigins = []string{"*"}
	}
	return cors.New(conf)
}
