package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"

	"github.com/gin-gonic/gin"
)

func GetRole(ctx *gin.Context) {
	role := middleware.GetRole(ctx)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetRoleResp(role)))
}

func GetRolePermissions(ctx *gin.Context) {
	role, ret := db.InitRoleRepo(db.DB).GetByID(middleware.GetRole(ctx).ID, db.GetOptions{
		Preloads: map[string]db.GetOptions{"Permissions": {}},
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	permissions := make([]gin.H, 0, len(role.Permissions))
	for _, perm := range role.Permissions {
		permissions = append(permissions, resp.GetPermissionResp(perm))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"permissions": permissions}))
}

func GetRoles(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	roles, count, ret := db.InitRoleRepo(db.DB).List(form.Limit, form.Offset)
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
	role, ret := db.InitRoleRepo(db.DB).Create(db.CreateRoleOptions{
		Name:        form.Name,
		Description: form.Description,
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
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
	if role.Default && form.Name != nil {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Model.Role.CannotUpdateDefault})
		return
	}
	ret := db.InitRoleRepo(db.DB).Update(role.ID, db.UpdateRoleOptions{
		Name:        form.Name,
		Description: form.Description,
	})
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func DeleteRole(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteRoleEventType)
	role := middleware.GetRole(ctx)
	if role.Default {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Model.Role.CannotDeleteDefault})
		return
	}
	ret := db.InitRoleRepo(db.DB).Delete(role.ID)
	if ret.OK {
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
	permission, ret := db.InitPermissionRepo(db.DB).GetByID(form.PermissionID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ret = db.AssignPermissionToRole(db.DB, permission, role)
	if ret.OK {
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
	permission, ret := db.InitPermissionRepo(db.DB).GetByID(form.PermissionID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ret = db.RevokePermissionFromRole(db.DB, permission, role)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}
