package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// CheckLogin 是否登录
func CheckLogin(ctx *gin.Context) {
	auth := strings.Fields(ctx.GetHeader("Authorization"))
	if len(auth) != 2 || auth[0] != "Bearer" {
		msg := "Unauthorized"
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	claims, err := utils.Parse(auth[1])
	if err != nil {
		msg := "Unauthorized"
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	if claims.Type == "admin" {
		admin, ok, msg := db.GetAdminByID(db.DB.WithContext(ctx), claims.UserID)
		if !ok {
			ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
			ctx.Abort()
			return
		}
		ctx.Set("Role", "admin")
		ctx.Set("Self", admin)
		ctx.Next()
	} else if claims.Type == "user" {
		user, ok, msg := db.GetUserByID(db.DB.WithContext(ctx), claims.UserID, false)
		if !ok {
			ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
			ctx.Abort()
			return
		}
		if magic := ctx.GetHeader("X-M"); magic != "" {
			tx := db.DB.WithContext(ctx).Begin()
			db.RecordDevice(tx, model.Device{UserID: user.ID, Magic: magic})
			if utils.EncryptMagic(magic) != claims.X {
				db.CreateCheat(tx, user.ID, 0, 0, model.MagicNotMatch, model.Suspect)
			}
			tx.Commit()
		}
		if user.Banned {
			ctx.JSONP(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
			ctx.Abort()
			return
		}
		ctx.Set("Role", "user")
		ctx.Set("Self", user)
		ctx.Next()
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
		ctx.Abort()
	}
}

// GetRole 获取角色， user 或 admin
func GetRole(ctx *gin.Context) string {
	if role, ok := ctx.Get("Role"); !ok {
		return ""
	} else {
		return role.(string)
	}
}

// GetSelfID 获取当前登录 admin 或 user 的ID
func GetSelfID(ctx *gin.Context) uint {
	var id uint
	switch GetRole(ctx) {
	case "admin":
		if self, ok := GetSelf(ctx).(model.Admin); ok {
			id = self.ID
		}
	case "user":
		if self, ok := GetSelf(ctx).(model.User); ok {
			id = self.ID
		}
	default:
		id = 0
	}
	return id
}

// GetSelf 获取当前登录 admin 或 user
func GetSelf(ctx *gin.Context) interface{} {
	if self, ok := ctx.Get("Self"); !ok {
		return nil
	} else {
		return self
	}
}

// CheckRole 检查角色
func CheckRole(t string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if GetRole(ctx) != t {
			ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// CheckCaptain 检查是否为队伍队长, 要求 uri 中必须包含 contestID, admin 路由不能使用
func CheckCaptain(ctx *gin.Context) {
	team := GetTeam(ctx)
	if team.CaptainID != GetSelfID(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
		ctx.Abort()
		return
	}
	ctx.Next()
}

// CheckVerified 检查邮箱是否已验证
func CheckVerified(ctx *gin.Context) {
	if self, ok := GetSelf(ctx).(model.User); GetRole(ctx) == "user" && ok && !self.Verified {
		ctx.JSON(http.StatusOK, gin.H{"msg": "UnverifiedEmail", "data": nil})
		ctx.Abort()
		return
	}
	ctx.Next()
}

// CheckBanned 检查队伍是否被封禁
func CheckBanned(ctx *gin.Context) {
	team := GetTeam(ctx)
	if team.Banned {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
		ctx.Abort()
		return
	}
	ctx.Next()
}
