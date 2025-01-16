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

func AdminChangePassword(ctx *gin.Context) {
	self, ok := ctx.Get("Self")
	trace := middleware.GetTraceID(ctx)
	if ok || self.(map[string]interface{})["Type"].(string) == "admin" {
		var form ChangePasswordForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "Bad Request"), "data": nil})
			return
		}
		ok, msg := db.ChangePasswordAdmin(ctx, self.(map[string]interface{})["ID"].(uint), form.OldPassword, form.NewPassword)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, "Success"), "data": nil})
		}
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"trace": trace, "msg": utils.M(ctx, "Forbidden"), "data": nil})
	}
}
