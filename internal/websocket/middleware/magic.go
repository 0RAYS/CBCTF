package middleware

import "github.com/gin-gonic/gin"

func SetMagic(ctx *gin.Context) {
	ctx.Set("Magic", ctx.Query("m"))
	ctx.Next()
}

func GetMagic(ctx *gin.Context) string {
	return ctx.GetString("Magic")
}
