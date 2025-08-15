package middleware

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

// GetSelf 获取当前登录 admin 或 user
func GetSelf(ctx *gin.Context) any {
	if self, ok := ctx.Get("Self"); !ok || self == nil {
		return nil
	} else {
		return self
	}
}

// GetSelfID 获取当前登录 admin 或 user 的ID
func GetSelfID(ctx *gin.Context) uint {
	var id uint
	if IsAdmin(ctx) {
		if self, ok := GetSelf(ctx).(model.Admin); ok {
			id = self.ID
		}
	} else {
		if self, ok := GetSelf(ctx).(model.User); ok {
			id = self.ID
		}
	}
	return id
}

func IsAdmin(ctx *gin.Context) bool {
	return ctx.GetBool("IsAdmin")
}
