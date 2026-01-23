package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUser(ctx *gin.Context) {
	var user model.User
	if middleware.IsAdmin(ctx) {
		user = middleware.GetUser(ctx)
	} else {
		user = middleware.GetSelf(ctx).(model.User)
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetUserResp(user, middleware.IsAdmin(ctx))))
}

func GetUsers(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	users, count, ret := db.InitUserRepo(db.DB).List(form.Limit, form.Offset)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, user := range users {
		data = append(data, resp.GetUserResp(user, true))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"count": count, "users": data}))
}

func CreateUser(ctx *gin.Context) {
	var form dto.CreateUserForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateUserEventType)
	user, ret := service.AdminCreateUser(db.DB, form)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetUserResp(user, true)))
}

func ChangePwd(ctx *gin.Context) {
	var form dto.ChangePasswordForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateUserEventType)
	ret := service.ChangeUserPwd(db.DB, middleware.GetSelf(ctx).(model.User), form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}

func UpdateUser(ctx *gin.Context) {
	var (
		user model.User
		ret  model.RetVal
	)
	if middleware.IsAdmin(ctx) {
		var form dto.UpdateUserForm
		if ret = form.Bind(ctx); !ret.OK {
			ctx.JSON(http.StatusOK, ret)
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UpdateUserEventType)
		user = middleware.GetUser(ctx)
		ret = service.UpdateUser(db.DB, user, form)
	} else {
		var form dto.UpdateSelfForm
		if ret = form.Bind(ctx); !ret.OK {
			ctx.JSON(http.StatusOK, ret)
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UpdateUserEventType)
		user = middleware.GetSelf(ctx).(model.User)
		ret = service.UpdateSelf(db.DB, user, form)
	}
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}

func DeleteUser(ctx *gin.Context) {
	var (
		tx  *gorm.DB
		ret model.RetVal
	)
	if !middleware.IsAdmin(ctx) {
		var form dto.DeleteSelfForm
		if ret = form.Bind(ctx); !ret.OK {
			ctx.JSON(http.StatusOK, ret)
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.DeleteUserEventType)
		tx = db.DB.Begin()
		ret = service.DeleteSelf(tx, middleware.GetSelf(ctx).(model.User), form)
	} else {
		ctx.Set(middleware.CTXEventTypeKey, model.DeleteUserEventType)
		tx = db.DB.Begin()
		ret = db.InitUserRepo(tx).Delete(middleware.GetUser(ctx).ID)
	}
	if !ret.OK {
		tx.Rollback()
	} else {
		tx.Commit()
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}
