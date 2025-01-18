package router

import (
	"CBCTF/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAdmin(ctx *gin.Context) {
	self, _ := ctx.Get("Self")
	admin, ok, msg := db.GetAdminByID(ctx, self.(map[string]interface{})["ID"].(uint))
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": msg, "data": nil})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": admin})
	}
}

func AdminChangePassword(ctx *gin.Context) {
	self, _ := ctx.Get("Self")
	var form ChangePasswordForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	ok, msg := db.ChangePasswordAdmin(ctx, self.(map[string]interface{})["ID"].(uint), form.OldPassword, form.NewPassword)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
	}
}
