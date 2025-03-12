package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func SetMagic(ctx *gin.Context) {
	magic := ctx.GetHeader("X-M")
	if strings.HasPrefix(ctx.Request.RequestURI, "/admin/") {
		ctx.Set("Magic", magic)
		ctx.Next()
		return
	}
	if !(strings.HasPrefix(ctx.Request.RequestURI, "/avatars/") && len(ctx.Request.RequestURI) == 45) && magic == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Magic", magic)
	ctx.Next()
}

func GetMagic(ctx *gin.Context) string {
	return ctx.GetString("Magic")
}
