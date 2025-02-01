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
		admin, ok, msg := db.GetAdminByID(ctx, claims.UserID)
		ctx.Set("Verified", true)
		ctx.Set("Hidden", false)
		if !ok {
			ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
			ctx.Abort()
			return
		}
		ctx.Set("SelfID", admin.ID)
		ctx.Set("Role", "admin")
		ctx.Set("Self", admin)
		ctx.Next()
		return
	} else if claims.Type == "user" {
		user, ok, msg := db.GetUserByID(ctx, claims.UserID, false)
		ctx.Set("Verified", user.Verified)
		ctx.Set("Hidden", user.Banned)
		if !ok {
			ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
			ctx.Abort()
			return
		}
		if user.Banned {
			ctx.JSONP(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
			ctx.Abort()
			return
		}
		ctx.Set("SelfID", user.ID)
		ctx.Set("Role", "user")
		ctx.Set("Self", user)
		ctx.Next()
		return
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
		ctx.Abort()
		return
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
	if selfID, ok := ctx.Get("SelfID"); !ok {
		return 0
	} else {
		return selfID.(uint)
	}
}

// GetSelf 获取当前登录 admin 或 user
func GetSelf(ctx *gin.Context) interface{} {
	if self, ok := ctx.Get("Self"); !ok {
		return nil
	} else {
		switch GetRole(ctx) {
		case "admin":
			return self.(model.Admin)
		case "user":
			return self.(model.User)
		default:
			return nil
		}
	}
}

// CheckRole 检查角色
func CheckRole(t string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if GetRole(ctx) != t {
			ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
			ctx.Abort()
		}
		ctx.Next()
	}
}

// CheckCaptain 检查是否为队伍队长, 要求 uri 中必须包含 contestID, admin 路由不能使用
func CheckCaptain(ctx *gin.Context) {
	team, ok, msg := db.GetTeamByUserID(ctx, GetSelfID(ctx), GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	if team.CaptainID != GetSelfID(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("TeamID", team.ID)
	ctx.Next()
}

// CheckVerified 检查邮箱是否已验证
func CheckVerified(ctx *gin.Context) {
	verified, ok := ctx.Get("Verified")
	if !ok || !verified.(bool) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "UnverifiedEmail", "data": nil})
		ctx.Abort()
	}
	ctx.Next()
}

func CheckBanned(ctx *gin.Context) {
	team, ok, msg := db.GetTeamByUserID(ctx, GetSelfID(ctx), GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	if team.Banned {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
		ctx.Abort()
		return
	}
	ctx.Next()
}
