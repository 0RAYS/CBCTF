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

func GetUser(ctx *gin.Context) {
	var user model.User
	includeCounts := middleware.IsFullAccess(ctx)
	if middleware.IsFullAccess(ctx) {
		user = middleware.GetUser(ctx)
	} else {
		user = middleware.GetSelf(ctx)
	}
	resp.JSON(ctx, model.SuccessRetVal(resp.GetUserResp(service.GetUserView(db.DB, user, includeCounts), includeCounts)))
}

func GetAccessibleRoutes(ctx *gin.Context) {
	userID := middleware.GetSelf(ctx).ID
	routes, ret := service.GetAccessibleRoutes(db.DB, userID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(routes))
}

func GetUsers(ctx *gin.Context) {
	var form dto.ListUsersForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	users, count, ret := service.ListUsers(db.DB, form)
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
	resp.JSON(ctx, model.SuccessRetVal(resp.GetUserResp(service.GetUserView(db.DB, user, true), true)))
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
	var ret model.RetVal
	var user model.User
	if !middleware.IsFullAccess(ctx) {
		var form dto.DeleteSelfForm
		if ret = dto.Bind(ctx, &form); !ret.OK {
			resp.JSON(ctx, ret)
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.DeleteUserEventType)
		user = middleware.GetSelf(ctx)
		ret = service.DeleteSelfWithTransaction(db.DB, user, form)
	} else {
		ctx.Set(middleware.CTXEventTypeKey, model.DeleteUserEventType)
		user = middleware.GetUser(ctx)
		ret = service.DeleteUserWithTransaction(db.DB, user)
	}
	if !ret.OK {
		resp.JSON(ctx, ret)
	} else {
		redis.DeleteUserRBAC(user.ID)
		ctx.Set(middleware.CTXEventSuccessKey, true)
		resp.JSON(ctx, ret)
	}
}
