package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CheckAuth 是否登录, 用户是否被 ban, 记录设备
func CheckAuth(ctx *gin.Context) {
	auth := strings.Fields(ctx.GetHeader("Authorization"))
	DB := db.DB.WithContext(ctx)
	if len(auth) != 2 || auth[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Unauthorized, "data": nil})
		return
	}
	claims, err := utils.ParseToken(auth[1])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Unauthorized, "data": nil})
		return
	}
	if claims.IsAdmin {
		admin, ok, msg := db.InitAdminRepo(DB).GetByID(claims.UserID)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		ctx.Set("IsAdmin", true)
		ctx.Set("Self", admin)
		ctx.Next()
	} else {
		user, ok, msg := db.InitUserRepo(DB).GetByID(claims.UserID)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		magic := GetMagic(ctx)
		if !utils.CompareMagic(magic, claims.X) {
			if utils.CompareMagic(model.OauthLoginType, claims.X) {
				token, err := utils.GenerateToken(user.ID, user.Name, false, magic)
				if err != nil {
					log.Logger.Warningf("Failed to generate token: %s", err)
					ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
					return
				}
				ctx.Writer.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
			} else {
				go db.InitCheatRepo(db.DB.WithContext(ctx.Copy())).Create(db.CreateCheatOptions{
					UserID:  sql.Null[uint]{V: user.ID, Valid: true},
					Magic:   magic,
					IP:      ctx.ClientIP(),
					Reason:  fmt.Sprintf(model.DifferentTokenMagic, magic, claims.X),
					Type:    model.Suspicious,
					Checked: false,
				})
				ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Unauthorized, "data": nil})
				return
			}
		}
		go db.InitDeviceRepo(db.DB.WithContext(ctx.Copy())).RecordDevice(user.ID, magic, ctx.ClientIP())
		if user.Banned {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Forbidden, "data": nil})
			return
		}
		ctx.Set("IsAdmin", false)
		ctx.Set("Self", user)
		ctx.Next()
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

// CheckRole 检查角色
func CheckRole(isAdmin bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if IsAdmin(ctx) != isAdmin {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Forbidden, "data": nil})
			return
		}
		ctx.Next()
	}
}

// CheckCaptain 检查是否为队伍队长, 要求 uri 中必须包含 contestID, admin 路由不能使用
func CheckCaptain(ctx *gin.Context) {
	team := GetTeam(ctx)
	if team.CaptainID != GetSelfID(ctx) {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Forbidden, "data": nil})
		return
	}
	ctx.Next()
}

// CheckVerified 检查邮箱是否已验证
func CheckVerified(ctx *gin.Context) {
	if self, ok := GetSelf(ctx).(model.User); !IsAdmin(ctx) && ok && !self.Verified {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.UnverifiedEmail, "data": nil})
		return
	}
	ctx.Next()
}

// CheckBanned 检查队伍是否被封禁
func CheckBanned(ctx *gin.Context) {
	team := GetTeam(ctx)
	if team.Banned {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Forbidden, "data": nil})
		return
	}
	ctx.Next()
}
