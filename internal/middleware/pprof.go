package middleware

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"net"

	"github.com/gin-gonic/gin"
)

func PProfWhitelist(ctx *gin.Context) {
	clientIP := ctx.ClientIP()
	ip := net.ParseIP(clientIP)
	for _, entry := range config.Env.Gin.PProf.Whitelist {
		if entry == clientIP {
			ctx.Next()
			return
		}
		_, cidr, err := net.ParseCIDR(entry)
		if err == nil && ip != nil && cidr.Contains(ip) {
			ctx.Next()
			return
		}
	}
	resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.Forbidden})
}
