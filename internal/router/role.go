package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRole(ctx *gin.Context) {
	roles, ret := db.InitRoleRepo(db.DB).GetByID(middleware.GetRole(ctx).ID, db.GetOptions{
		Preloads: map[string]db.GetOptions{"Permissions": {}},
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetRoleResp(roles)))
}

func GetRoles(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	roles, count, ret := db.InitRoleRepo(db.DB).List(form.Limit, form.Offset)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, role := range roles {
		data = append(data, resp.GetRoleResp(role))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"count": count, "roles": data}))
}

func CreateRole(ctx *gin.Context) {
	var form dto.CreateRoleForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateRoleEventType)
	role, ret := db.InitRoleRepo(db.DB).Create(db.CreateRoleOptions{
		Name:        form.Name,
		Description: form.Description,
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetRoleResp(role)))
}

func UpdateRole(ctx *gin.Context) {
	var form dto.UpdateRoleForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateRoleEventType)
	role := middleware.GetRole(ctx)
	if role.Default && form.Name != nil {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Request.Forbidden})
		return
	}
	ret := db.InitRoleRepo(db.DB).Update(role.ID, db.UpdateRoleOptions{
		Name:        form.Name,
		Description: form.Description,
	})
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}

func DeleteRole(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteRoleEventType)
	role := middleware.GetRole(ctx)
	if role.Default {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Request.Forbidden})
		return
	}
	ret := db.InitRoleRepo(db.DB).Delete(role.ID)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}

func AssignPermission(ctx *gin.Context) {
	var form dto.AssignPermissionForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.AssignPermissionEventType)
	role := middleware.GetRole(ctx)
	permission, ret := db.InitPermissionRepo(db.DB).GetByID(form.PermissionID)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ret = db.AssignPermissionToRole(db.DB, permission, role)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}

func RevokePermission(ctx *gin.Context) {
	var form dto.AssignPermissionForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.RevokePermissionEventType)
	role := middleware.GetRole(ctx)
	permission, ret := db.InitPermissionRepo(db.DB).GetByID(form.PermissionID)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ret = db.RevokePermissionFromRole(db.DB, permission, role)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}
