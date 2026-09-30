package middleware

import (
	"net/url"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"CBCTF/internal/config"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	var origins []string
	for _, origin := range config.Env.Gin.Origins {
		u, err := url.Parse(origin)
		if err != nil {
			continue
		}
		origins = append(origins, strings.Trim(u.String(), "/"))
	}
	conf := cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "Connection", "Upgrade",
		},
		ExposeHeaders: []string{
			"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control",
			"Content-Language", "Content-Type", "Authorization", "File",
		},
		AllowCredentials: true,
	}
	return cors.New(conf)
}
