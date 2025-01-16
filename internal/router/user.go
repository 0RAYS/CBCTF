package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUser(ctx *gin.Context) {
	self, ok := ctx.Get("Self")
	trace := middleware.GetTraceID(ctx)
	if ok || self.(map[string]interface{})["Type"].(string) == "user" {
		user, ok, msg := db.GetUserByID(ctx, self.(map[string]interface{})["ID"].(uint), true)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, "Success"), "data": user})
		}
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"trace": trace, "msg": utils.M(ctx, "Forbidden"), "data": nil})
	}
}
