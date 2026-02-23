package middleware

import (
	"github.com/gin-gonic/gin"
)

func SetMagic(ctx *gin.Context) {
	_, magic := parseWSToken(ctx.Request.Header.Get("Sec-Websocket-Protocol"))
	ctx.Set("Magic", magic)
	ctx.Next()
}

func GetMagic(ctx *gin.Context) string {
	return ctx.GetString("Magic")
}
