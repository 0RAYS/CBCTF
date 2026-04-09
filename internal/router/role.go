package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"

	"github.com/gin-gonic/gin"
)

func GetRole(ctx *gin.Context) {
	role := middleware.GetRole(ctx)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetRoleResp(role)))
}

func GetRolePermissions(ctx *gin.Context) {
	permissions, ret := service.GetRolePermissions(db.DB, middleware.GetRole(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(permissions))
	for _, perm := range permissions {
		data = append(data, resp.GetPermissionResp(perm))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"permissions": data}))
}

func GetRoles(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	roles, count, ret := service.ListRoles(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, role := range roles {
		data = append(data, resp.GetRoleResp(role))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "roles": data}))
}

func CreateRole(ctx *gin.Context) {
	var form dto.CreateRoleForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateRoleEventType)
	role, ret := service.CreateRole(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	redis.DeleteRBAC()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetRoleResp(role)))
}

func UpdateRole(ctx *gin.Context) {
	var form dto.UpdateRoleForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateRoleEventType)
	role := middleware.GetRole(ctx)
	ret := service.UpdateRole(db.DB, role, form)
	if ret.OK {
		redis.DeleteRBAC()
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func DeleteRole(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteRoleEventType)
	role := middleware.GetRole(ctx)
	ret := service.DeleteRole(db.DB, role)
	if ret.OK {
		redis.DeleteRBAC()
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func AssignPermission(ctx *gin.Context) {
	var form dto.AssignPermissionForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.AssignPermissionEventType)
	role := middleware.GetRole(ctx)
	ret := service.AssignPermission(db.DB, role, form)
	if ret.OK {
		redis.DeleteRBAC()
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func RevokePermission(ctx *gin.Context) {
	var form dto.AssignPermissionForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.RevokePermissionEventType)
	role := middleware.GetRole(ctx)
	ret := service.RevokePermission(db.DB, role, form)
	if ret.OK {
		redis.DeleteRBAC()
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}
