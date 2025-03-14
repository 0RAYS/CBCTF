package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAdmin(ctx *gin.Context) {
	admin := middleware.GetSelf(ctx).(model.Admin)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": &admin})
}

func GetAdmins(ctx *gin.Context) {
	admins, count, ok, msg := db.GetAdmins(db.DB.WithContext(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "admins": &admins}})
}

func CreateAdmin(ctx *gin.Context) {
	var form f.CreateAdminForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	admin, ok, msg := db.CreateAdmin(tx, form.Name, form.Password, form.Email)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": &admin})
}

func AdminChangePassword(ctx *gin.Context) {
	var form f.ChangePasswordForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.ChangePasswordAdmin(tx, middleware.GetSelf(ctx).(model.Admin), form.OldPassword, form.NewPassword)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateAdmin(ctx *gin.Context) {
	var form f.UpdateAdminForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	admin := middleware.GetSelf(ctx).(model.Admin)
	data := utils.Form2Map(form)
	tx := db.DB.WithContext(ctx).Begin()
	// 在预期的想法中, admin 的邮箱似乎没有什么用, 先保留
	if email, ok := data["email"]; ok && email.(string) != admin.Email {
		if !db.IsUniqueEmail(tx, data["email"].(string)) {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": "EmailExists", "data": nil})
			return
		}
		if ok, msg := db.UpdateAdmin(tx, admin.ID, map[string]interface{}{"verified": false}); !ok {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
	}
	if name, ok := data["name"]; ok && name.(string) != admin.Name && !db.IsUniqueName(tx, name.(string), model.Admin{}) {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": "UserNameExists", "data": nil})
		return
	}
	ok, msg := db.UpdateAdmin(tx, admin.ID, data)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
