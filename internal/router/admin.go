package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAdmin(ctx *gin.Context) {
	admin := middleware.GetSelf(ctx).(model.Admin)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": resp.GetAdminResp(admin)})
}

func AdminChangePassword(ctx *gin.Context) {
	var form f.ChangePasswordForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateAdminEventType)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.ChangeAdminPassword(tx, middleware.GetSelf(ctx).(model.Admin), form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateAdmin(ctx *gin.Context) {
	var form f.UpdateAdminForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateAdminEventType)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.UpdateAdmin(tx, middleware.GetSelf(ctx).(model.Admin), form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func CreateAdmin(ctx *gin.Context) {
	var form f.CreateAdminForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateAdminEventType)
	tx := db.DB.WithContext(ctx).Begin()
	admin, ok, msg := service.CreateAdmin(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.GetAdminResp(admin)})
}
