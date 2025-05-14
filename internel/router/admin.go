package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAdmin(ctx *gin.Context) {
	admin := middleware.GetSelf(ctx).(model.Admin)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": resp.GetAdminResp(admin)})
}

func AdminChangePassword(ctx *gin.Context) {
	var form f.ChangePasswordForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.ChangeAdminPassword(tx, middleware.GetSelf(ctx).(model.Admin), form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateAdmin(ctx *gin.Context) {
	var form f.UpdateAdminForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.UpdateAdmin(tx, middleware.GetSelf(ctx).(model.Admin), form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func CreateAdmin(ctx *gin.Context) {
	var form f.CreateAdminForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	admin, ok, msg := service.CreateAdmin(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.GetAdminResp(admin)})
}
