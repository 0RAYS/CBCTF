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

func GetGroup(ctx *gin.Context) {
	group := middleware.GetGroup(ctx)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetGroupResp(service.GetGroupView(db.DB, group))))
}

func GetGroupUsers(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	group := middleware.GetGroup(ctx)
	users, count, ret := service.ListGroupUsers(db.DB, group, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(users))
	for _, user := range users {
		data = append(data, resp.GetUserResp(user, true))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "users": data}))
}

func GetGroupAvailableUsers(ctx *gin.Context) {
	var form dto.ListUsersForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	group := middleware.GetGroup(ctx)
	users, count, ret := service.ListUsersNotInGroup(db.DB, group, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(users))
	for _, user := range users {
		data = append(data, resp.GetUserResp(user, true))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "users": data}))
}

func GetGroups(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	groups, count, ret := service.ListGroups(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, group := range groups {
		data = append(data, resp.GetGroupResp(group))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "groups": data}))
}

func CreateGroup(ctx *gin.Context) {
	var form dto.CreateGroupForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateGroupEventType)
	group, ret := service.CreateGroup(db.DB, form)
	if !ret.OK {
		redis.DeleteRBAC()
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetGroupResp(service.GetGroupView(db.DB, group))))
}

func UpdateGroup(ctx *gin.Context) {
	var form dto.UpdateGroupForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateGroupEventType)
	group := middleware.GetGroup(ctx)
	ret := service.UpdateGroup(db.DB, group, form)
	if ret.OK {
		redis.DeleteRBAC()
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func DeleteGroup(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteGroupEventType)
	group := middleware.GetGroup(ctx)
	ret := service.DeleteGroup(db.DB, group)
	if ret.OK {
		redis.DeleteRBAC()
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func AssignUserToGroup(ctx *gin.Context) {
	var form dto.AssignUserGroupForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.AssignUserGroupEventType)
	group := middleware.GetGroup(ctx)
	user, ret := service.AssignUserToGroup(db.DB, group, form)
	if ret.OK {
		redis.DeleteUserRBAC(user.ID)
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func RemoveUserFromGroup(ctx *gin.Context) {
	var form dto.AssignUserGroupForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.RemoveUserGroupEventType)
	group := middleware.GetGroup(ctx)
	user, ret := service.RemoveUserFromGroup(db.DB, group, form)
	if ret.OK {
		redis.DeleteUserRBAC(user.ID)
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}
