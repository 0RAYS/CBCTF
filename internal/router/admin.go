package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
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
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetAdminResp(admin)))
}

func AdminChangePassword(ctx *gin.Context) {
	var form f.ChangePasswordForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateAdminEventType)
	ret := service.ChangeAdminPassword(db.DB, middleware.GetSelf(ctx).(model.Admin), form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}

func UpdateAdmin(ctx *gin.Context) {
	var form f.UpdateAdminForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateAdminEventType)
	ret := service.UpdateAdmin(db.DB, middleware.GetSelf(ctx).(model.Admin), form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}

func CreateAdmin(ctx *gin.Context) {
	var form f.CreateAdminForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateAdminEventType)
	admin, ret := service.CreateAdmin(db.DB, form)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetAdminResp(admin)))
}

func GetLogs(ctx *gin.Context) {
	var form f.GetLogsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data, _ := redis.GetLogs(int64(form.Offset), int64(form.Offset+form.Limit))
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}
