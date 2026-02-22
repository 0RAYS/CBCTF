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

func GetGroup(ctx *gin.Context) {
	group := middleware.GetGroup(ctx)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetGroupResp(group)))
}

func GetGroupUsers(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	group := middleware.GetGroup(ctx)
	users, count, ret := db.GetGroupUsers(db.DB, group, form.Limit, form.Offset)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0, len(users))
	for _, user := range users {
		data = append(data, resp.GetUserResp(user, true))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"count": count, "users": data}))
}

func GetGroups(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	groups, count, ret := db.InitGroupRepo(db.DB).List(form.Limit, form.Offset)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, group := range groups {
		data = append(data, resp.GetGroupResp(group))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"count": count, "groups": data}))
}

func CreateGroup(ctx *gin.Context) {
	var form dto.CreateGroupForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateGroupEventType)
	group, ret := db.InitGroupRepo(db.DB).Create(db.CreateGroupOptions{
		RoleID:      form.RoleID,
		Name:        form.Name,
		Description: form.Description,
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetGroupResp(group)))
}

func UpdateGroup(ctx *gin.Context) {
	var form dto.UpdateGroupForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateGroupEventType)
	group := middleware.GetGroup(ctx)
	if group.Default && form.Name != nil {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Request.Forbidden})
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
	ctx.JSON(http.StatusOK, ret)
}

func DeleteGroup(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteGroupEventType)
	group := middleware.GetGroup(ctx)
	if group.Default {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Request.Forbidden})
		return
	}
	ret := db.InitGroupRepo(db.DB).Delete(group.ID)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}

func AssignUserToGroup(ctx *gin.Context) {
	var form dto.AssignUserGroupForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.AssignUserGroupEventType)
	group := middleware.GetGroup(ctx)
	user, ret := db.InitUserRepo(db.DB).GetByID(form.UserID)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ret = db.AppendUserToGroup(db.DB, user, group)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}

func RemoveUserFromGroup(ctx *gin.Context) {
	var form dto.AssignUserGroupForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.RemoveUserGroupEventType)
	group := middleware.GetGroup(ctx)
	user, ret := db.InitUserRepo(db.DB).GetByID(form.UserID)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ret = db.DeleteUserFromGroup(db.DB, user, group)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}
