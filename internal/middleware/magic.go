package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetMagic(strict bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		magic := ctx.GetHeader("X-M")
		if strict && magic == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			ctx.Abort()
			return
		}
		ctx.Set("Magic", magic)
		ctx.Next()
	}
}

func GetMagic(ctx *gin.Context) string {
	return ctx.GetString("Magic")
}
