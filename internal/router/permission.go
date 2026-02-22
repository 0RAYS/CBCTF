package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPermissions(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	permissions, count, ret := db.InitPermissionRepo(db.DB).List(form.Limit, form.Offset)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, perm := range permissions {
		data = append(data, resp.GetPermissionResp(perm))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"count": count, "permissions": data}))
}

func UpdatePermission(ctx *gin.Context) {
	var form dto.UpdatePermissionForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
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
	ctx.JSON(http.StatusOK, ret)
}
