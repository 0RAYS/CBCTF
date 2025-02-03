package router

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAdmin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": middleware.GetSelf(ctx).(model.Admin)})
}

func GetAdmins(ctx *gin.Context) {
	admins, count, ok, msg := db.GetAdmins(ctx)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "admins": admins}})
}

func CreateAdmin(ctx *gin.Context) {
	var form constants.CreateAdminForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	admin, ok, msg := db.CreateAdmin(ctx, form.Name, form.Password, form.Email)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": admin})
}

func AdminChangePassword(ctx *gin.Context) {
	var form constants.ChangePasswordForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	_, msg := db.ChangePasswordAdmin(ctx, middleware.GetSelf(ctx).(model.Admin), form.OldPassword, form.NewPassword)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateAdmin(ctx *gin.Context) {
	var form constants.UpdateAdminForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	admin := middleware.GetSelf(ctx).(model.Admin)
	data := utils.Form2Map(form)
	// 在预期的想法中，admin 的邮箱似乎没有什么用，先保留
	if email, ok := data["email"]; ok && email.(string) != admin.Email {
		if !db.IsUniqueEmail(data["email"].(string)) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "EmailExists", "data": nil})
			return
		}
		db.UpdateUser(ctx, admin.ID, map[string]interface{}{"verified": false})
	}
	if name, ok := data["name"]; ok && name.(string) != admin.Name && !db.IsUniqueName(name.(string), model.Admin{}) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "UserNameExists", "data": nil})
		return
	}
	_, msg := db.UpdateAdmin(ctx, admin.ID, data)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
