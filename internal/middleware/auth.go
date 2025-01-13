package middleware

import (
	"RayWar/internal/db"
	"RayWar/internal/log"
	"RayWar/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// CheckLogin 是否登录
func CheckLogin(ctx *gin.Context) {
	auth := strings.Fields(ctx.GetHeader("Authorization"))
	trace := GetTraceID(ctx)
	if len(auth) != 2 || auth[0] != "Bearer" {
		msg := "Unauthorized"
		log.Logger.Debugf("| %s | Unauthorized", trace)
		ctx.JSON(http.StatusUnauthorized, gin.H{"trace": trace, "msg": utils.M(ctx, msg)})
		ctx.Abort()
		return
	}
	claims, err := utils.Parse(auth[1])
	if err != nil {
		msg := "Unauthorized"
		log.Logger.Debugf("| %s | Unauthorized", trace)
		ctx.JSON(http.StatusUnauthorized, gin.H{"trace": trace, "msg": utils.M(ctx, msg)})
		ctx.Abort()
		return
	}
	user, ok, msg := db.GetUserByID(claims.UserID)
	if !ok {
		log.Logger.Infof("| %s | %s", trace, msg)
		ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg)})
		ctx.Abort()
		return
	}
	if user.Banned && user.Type != "admin" {
		log.Logger.Infof("| %s | %s", trace, "Forbidden")
		ctx.JSONP(http.StatusForbidden, gin.H{"trace": trace, "msg": utils.M(ctx, "Forbidden")})
		ctx.Abort()
		return
	}
	ctx.Set("self", user)
	ctx.Next()
}

// CheckAdmin 是否为Admin
func CheckAdmin(ctx *gin.Context) {
	trace := GetTraceID(ctx)
	self := GetSelf(ctx)
	if self.Type != "admin" {
		log.Logger.Debugf("| %s | Forbidden", trace)
		ctx.JSON(http.StatusForbidden, gin.H{"trace": trace, "msg": utils.M(ctx, "Forbidden")})
		ctx.Abort()
		return
	}
	ctx.Next()
}
