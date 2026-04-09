package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"

	"github.com/gin-gonic/gin"
)

func GetPermissions(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	permissions, count, ret := service.ListPermissions(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, perm := range permissions {
		data = append(data, resp.GetPermissionResp(perm))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "permissions": data}))
}

func UpdatePermission(ctx *gin.Context) {
	var form dto.UpdatePermissionForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdatePermissionEventType)
	permission := middleware.GetPermission(ctx)
	ret := service.UpdatePermission(db.DB, permission, form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}
