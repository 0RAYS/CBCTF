package middleware

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

// GetSelf 获取当前登录 admin 或 user
func GetSelf(ctx *gin.Context) model.User {
	self, ok := ctx.Get("Self")
	if !ok || self == nil {
		return model.User{}
	}
	return self.(model.User)
}
