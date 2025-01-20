package router

import (
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

func AdminChangePassword(ctx *gin.Context) {
	var form ChangePasswordForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	_, msg := db.ChangePasswordAdmin(ctx, middleware.GetSelfID(ctx), form.OldPassword, form.NewPassword)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateAdmin(ctx *gin.Context) {
	var form UpdateAdminForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	admin := middleware.GetSelf(ctx).(model.Admin)
	data := utils.Form2Map(form)
	// 在预期的想法中，admin 的邮箱似乎没有什么用，先保留
	if admin.Email != data["email"].(string) {
		if !db.IsUniqueEmail(data["email"].(string)) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "EmailExists", "data": nil})
			return
		}
		db.UpdateUser(ctx, admin.ID, map[string]interface{}{"verified": false})
	}
	if admin.Name != data["name"].(string) && !db.IsUniqueName(data["name"].(string), model.Admin{}) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "UserNameExists", "data": nil})
		return
	}
	_, msg := db.UpdateAdmin(ctx, admin.ID, data)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
