package middleware

import (
	"CBCTF/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

// SetMagic 保存设备 ID 至上下文
func SetMagic(ctx *gin.Context) {
	magic := ctx.GetHeader("X-M")
	path := ctx.Request.URL.Path
	for _, pattern := range config.Env.Gin.Magic.Whitelist {
		rgx := regexp.MustCompile(pattern)
		if rgx.MatchString(path) {
			ctx.Set("Magic", magic)
			ctx.Next()
			return
		}
	}
	if magic == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Magic", magic)
	ctx.Next()
}

// GetMagic 获取设备 ID
func GetMagic(ctx *gin.Context) string {
	return ctx.GetString("Magic")
}
