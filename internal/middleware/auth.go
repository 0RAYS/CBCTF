package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

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
		ctx.Set("Self", map[string]interface{}{"ID": admin.ID, "Type": "admin"})
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
		ctx.Set("Self", map[string]interface{}{"ID": user.ID, "Type": "user"})
		ctx.Next()
		return
	}
}

func CheckType(t string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		self, ok := ctx.Get("Self")
		if !ok || self.(map[string]interface{})["Type"].(string) != t {
			ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
			ctx.Abort()
		}
		ctx.Next()
	}
}

func CheckVerified(ctx *gin.Context) {
	verified, ok := ctx.Get("Verified")
	if !ok || !verified.(bool) {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "UnverifiedEmail", "data": nil})
		ctx.Abort()
	}
	ctx.Next()
}
