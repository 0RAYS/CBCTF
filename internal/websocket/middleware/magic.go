package middleware

import (
	"github.com/gin-gonic/gin"
)

func SetMagic(ctx *gin.Context) {
	protocols := parseWSProtocols(ctx.Request.Header.Get("Sec-Websocket-Protocol"))
	magic := protocols["Magic"]
	if magic == "" {
		// fallback: legacy query param
		magic = ctx.Query("m")
	}
	ctx.Set("Magic", magic)
	ctx.Next()
}

func GetMagic(ctx *gin.Context) string {
	return ctx.GetString("Magic")
}
