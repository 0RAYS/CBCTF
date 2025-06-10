package middleware

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"CBCTF/internel/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// CheckAuth 是否登录, 用户是否被 ban, 记录设备
func CheckAuth(ctx *gin.Context) {
	auth := strings.Fields(ctx.GetHeader("Authorization"))
	DB := db.DB.WithContext(ctx)
	if len(auth) != 2 || auth[0] != "Bearer" {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Unauthorized, "data": nil})
		ctx.Abort()
		return
	}
	claims, err := utils.Parse(auth[1])
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Unauthorized, "data": nil})
		ctx.Abort()
		return
	}
	if claims.Type == "admin" {
		admin, ok, msg := db.InitAdminRepo(DB).GetByID(claims.UserID, "all")
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			ctx.Abort()
			return
		}
		ctx.Set("Role", "admin")
		ctx.Set("Self", admin)
		ctx.Next()
	} else if claims.Type == "user" {
		user, ok, msg := db.InitUserRepo(DB).GetByID(claims.UserID, "all")
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			ctx.Abort()
			return
		}
		service.RecordDevice(DB, user.ID, GetMagic(ctx))
		if !utils.CompareMagic(GetMagic(ctx), claims.X) {
			db.InitCheatRepo(db.DB.WithContext(ctx)).Create(db.CreateCheatOptions{
				UserID:     &user.ID,
				Magic:      GetMagic(ctx),
				IP:         ctx.ClientIP(),
				Reason:     fmt.Sprintf("Device magic %s is different from token magic %s", GetMagic(ctx), claims.X),
				Type:       model.Suspicious,
				Checked:    false,
				References: nil,
			})
		}
		if user.Banned {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Forbidden, "data": nil})
			ctx.Abort()
			return
		}
		ctx.Set("Role", "user")
		ctx.Set("Self", user)
		ctx.Next()
	} else {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Forbidden, "data": nil})
		ctx.Abort()
	}
}

// GetRole 获取角色,  user 或 admin
func GetRole(ctx *gin.Context) string {
	if role, ok := ctx.Get("Role"); !ok || role == nil {
		return ""
	} else {
		return role.(string)
	}
}

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

// CheckRole 检查角色
func CheckRole(t string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if GetRole(ctx) != t {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Forbidden, "data": nil})
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
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Forbidden, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Next()
}

// CheckVerified 检查邮箱是否已验证
func CheckVerified(ctx *gin.Context) {
	if self, ok := GetSelf(ctx).(model.User); GetRole(ctx) == "user" && ok && !self.Verified {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnverifiedEmail, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Next()
}

// CheckBanned 检查队伍是否被封禁
func CheckBanned(ctx *gin.Context) {
	team := GetTeam(ctx)
	if team.Banned {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Forbidden, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Next()
}
