package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAdmin(ctx *gin.Context) {
	self, ok := ctx.Get("Self")
	trace := middleware.GetTraceID(ctx)
	if ok || self.(map[string]interface{})["Type"].(string) == "admin" {
		admin, ok, msg := db.GetAdminByID(ctx, self.(map[string]interface{})["ID"].(uint))
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, "Success"), "data": admin})
		}
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"trace": trace, "msg": utils.M(ctx, "Forbidden"), "data": nil})
	}
}
