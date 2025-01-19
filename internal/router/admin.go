package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
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

func UpdateAdmin(ctx *gin.Context) {
	self, _ := ctx.Get("Self")
	var form UpdateAdminForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	admin, ok, msg := db.GetAdminByID(ctx, self.(map[string]interface{})["ID"].(uint))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
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
	ok, msg = db.UpdateAdmin(ctx, admin.ID, data)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
	}
}

func UpdateUser(ctx *gin.Context) {
	var form UpdateUserForm
	type userIDUri struct {
		UserID uint `uri:"userID" binding:"required"`
	}
	var uri userIDUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	user, ok, msg := db.GetUserByID(ctx, uri.UserID, false)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := utils.Form2Map(form)
	if user.Email != data["email"].(string) {
		if !db.IsUniqueEmail(data["email"].(string)) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "EmailExists", "data": nil})
			return
		}
		db.UpdateUser(ctx, user.ID, map[string]interface{}{"verified": false})
	}
	if user.Name != data["name"].(string) && !db.IsUniqueName(data["name"].(string), model.User{}) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "UserNameExists", "data": nil})
		return
	}
	if data["password"] != "" {
		data["password"] = utils.HashPassword(data["password"].(string))
	} else {
		data["password"] = user.Password
	}
	ok, msg = db.UpdateUser(ctx, user.ID, data)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
	}
}
