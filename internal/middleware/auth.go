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
	trace := GetTraceID(ctx)
	if len(auth) != 2 || auth[0] != "Bearer" {
		msg := "Unauthorized"
		ctx.JSON(http.StatusUnauthorized, gin.H{"trace": trace, "msg": msg})
		ctx.Abort()
		return
	}
	claims, err := utils.Parse(auth[1])
	if err != nil {
		msg := "Unauthorized"
		ctx.JSON(http.StatusUnauthorized, gin.H{"trace": trace, "msg": msg})
		ctx.Abort()
		return
	}
	if claims.Type == "admin" {
		admin, ok, msg := db.GetAdminByID(ctx, claims.UserID)
		if !ok {
			ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": msg})
			ctx.Abort()
			return
		}
		ctx.Set("Self", map[string]interface{}{"ID": admin.ID, "Type": "admin"})
		ctx.Next()
		return
	} else if claims.Type == "user" {
		user, ok, msg := db.GetUserByID(ctx, claims.UserID, false)
		if !ok {
			ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": msg})
			ctx.Abort()
			return
		}
		if user.Banned {
			ctx.JSONP(http.StatusForbidden, gin.H{"trace": trace, "msg": "Forbidden"})
			ctx.Abort()
			return
		}
		ctx.Set("Self", map[string]interface{}{"ID": user.ID, "Type": "user"})
		ctx.Next()
		return
	}
}

func CheckAdmin(ctx *gin.Context) {
	trace := GetTraceID(ctx)
	self, ok := ctx.Get("Self")
	if !ok || self.(map[string]interface{})["Type"].(string) != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"trace": trace, "msg": "Forbidden"})
		ctx.Abort()
		return
	}
	ctx.Next()
}
