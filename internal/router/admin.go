package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
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
	ok, msg := service.ChangeAdminPassword(db.DB, middleware.GetSelf(ctx).(model.Admin), form)
	if ok {
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
	ok, msg := service.UpdateAdmin(db.DB, middleware.GetSelf(ctx).(model.Admin), form)
	if ok {
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
	admin, ok, msg := service.CreateAdmin(db.DB, form)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.GetAdminResp(admin)})
}

func GetLogs(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data, _, msg := redis.GetLogs(int64(form.Offset), int64(form.Limit))
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": data})
}
