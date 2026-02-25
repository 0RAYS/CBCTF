package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"

	"github.com/gin-gonic/gin"
)

func GetPermissions(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	permissions, count, ret := db.InitPermissionRepo(db.DB).List(form.Limit, form.Offset)
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
	ret := db.InitPermissionRepo(db.DB).Update(permission.ID, db.UpdatePermissionOptions{
		Description: form.Description,
	})
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}
