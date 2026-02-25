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

func GetGroup(ctx *gin.Context) {
	group := middleware.GetGroup(ctx)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetGroupResp(group)))
}

func GetGroupUsers(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	group := middleware.GetGroup(ctx)
	users, count, ret := db.InitUserRepo(db.DB).GetByGroupID(group.ID, form.Limit, form.Offset)
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
	groups, count, ret := db.InitGroupRepo(db.DB).List(form.Limit, form.Offset)
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
	group, ret := db.InitGroupRepo(db.DB).Create(db.CreateGroupOptions{
		RoleID:      form.RoleID,
		Name:        form.Name,
		Description: form.Description,
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetGroupResp(group)))
}

func UpdateGroup(ctx *gin.Context) {
	var form dto.UpdateGroupForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateGroupEventType)
	group := middleware.GetGroup(ctx)
	if group.Default && form.Name != nil {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Model.Group.CannotUpdateDefault})
		return
	}
	ret := db.InitGroupRepo(db.DB).Update(group.ID, db.UpdateGroupOptions{
		RoleID:      form.RoleID,
		Name:        form.Name,
		Description: form.Description,
	})
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func DeleteGroup(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteGroupEventType)
	group := middleware.GetGroup(ctx)
	if group.Default {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Model.Group.CannotDeleteDefault})
		return
	}
	ret := db.InitGroupRepo(db.DB).Delete(group.ID)
	if ret.OK {
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
	user, ret := db.InitUserRepo(db.DB).GetByID(form.UserID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ret = db.AppendUserToGroup(db.DB, user, group)
	if ret.OK {
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
	user, ret := db.InitUserRepo(db.DB).GetByID(form.UserID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ret = db.DeleteUserFromGroup(db.DB, user, group)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}
