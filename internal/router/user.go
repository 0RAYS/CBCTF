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
	"gorm.io/gorm"
)

func GetUser(ctx *gin.Context) {
	var user model.User
	if middleware.IsFullAccess(ctx) {
		user = middleware.GetUser(ctx)
	} else {
		user = middleware.GetSelf(ctx)
	}
	resp.JSON(ctx, model.SuccessRetVal(resp.GetUserResp(user, middleware.IsFullAccess(ctx))))
}

func GetAccessibleRoutes(ctx *gin.Context) {
	userID := middleware.GetSelf(ctx).ID
	permNames, ret := db.InitPermissionRepo(db.DB).GetUserPermissions(userID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}

	permSet := make(map[string]struct{}, len(permNames))
	for _, name := range permNames {
		permSet[name] = struct{}{}
	}

	routes := make([]string, 0)
	for route, perm := range model.RoutePermissions {
		if _, ok := permSet[perm]; ok {
			routes = append(routes, route)
		}
	}

	resp.JSON(ctx, model.SuccessRetVal(routes))
}

func GetUsers(ctx *gin.Context) {
	var form dto.ListUsersForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	options := db.GetOptions{Search: make(map[string]string)}
	if form.Name != "" {
		options.Search["name"] = form.Name
	}
	if form.Email != "" {
		options.Search["email"] = form.Email
	}
	if form.Description != "" {
		options.Search["description"] = form.Description
	}
	users, count, ret := db.InitUserRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, user := range users {
		data = append(data, resp.GetUserResp(user, true))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "users": data}))
}

func CreateUser(ctx *gin.Context) {
	var form dto.CreateUserForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateUserEventType)
	user, ret := service.AdminCreateUser(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetUserResp(user, true)))
}

func ChangePwd(ctx *gin.Context) {
	var form dto.ChangePasswordForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateUserEventType)
	ret := service.ChangeUserPwd(db.DB, middleware.GetSelf(ctx), form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func UpdateUser(ctx *gin.Context) {
	var (
		user model.User
		ret  model.RetVal
	)
	if middleware.IsFullAccess(ctx) {
		var form dto.UpdateUserForm
		if ret = dto.Bind(ctx, &form); !ret.OK {
			resp.JSON(ctx, ret)
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UpdateUserEventType)
		user = middleware.GetUser(ctx)
		ret = service.UpdateUser(db.DB, user, form)
	} else {
		var form dto.UpdateSelfForm
		if ret = dto.Bind(ctx, &form); !ret.OK {
			resp.JSON(ctx, ret)
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UpdateUserEventType)
		user = middleware.GetSelf(ctx)
		ret = service.UpdateSelf(db.DB, user, form)
	}
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func DeleteUser(ctx *gin.Context) {
	var (
		tx  *gorm.DB
		ret model.RetVal
	)
	var user model.User
	if !middleware.IsFullAccess(ctx) {
		var form dto.DeleteSelfForm
		if ret = dto.Bind(ctx, &form); !ret.OK {
			resp.JSON(ctx, ret)
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.DeleteUserEventType)
		user = middleware.GetSelf(ctx)
		tx = db.DB.Begin()
		ret = service.DeleteSelf(tx, user, form)
	} else {
		ctx.Set(middleware.CTXEventTypeKey, model.DeleteUserEventType)
		user = middleware.GetUser(ctx)
		tx = db.DB.Begin()
		ret = service.DeleteUser(tx, user)
	}
	if !ret.OK {
		tx.Rollback()
	} else {
		redis.DeleteUserRBAC(user.ID)
		tx.Commit()
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}
