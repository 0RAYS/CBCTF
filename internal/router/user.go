package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUser(ctx *gin.Context) {
	var user model.User
	if middleware.IsAdmin(ctx) {
		user = middleware.GetUser(ctx)
	} else {
		user = middleware.GetSelf(ctx).(model.User)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": resp.GetUserResp(user, middleware.IsAdmin(ctx))})
}

func GetUsers(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	users, count, ok, msg := db.InitUserRepo(db.DB).List(form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, user := range users {
		data = append(data, resp.GetUserResp(user, true))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "users": data}})
}

func CreateUser(ctx *gin.Context) {
	var form f.CreateUserForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateUserEventType)
	user, ok, msg := service.AdminCreateUser(db.DB, form)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.GetUserResp(user, true)})
}

func ChangePwd(ctx *gin.Context) {
	var form f.ChangePasswordForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateUserEventType)
	ok, msg := service.ChangeUserPwd(db.DB, middleware.GetSelf(ctx).(model.User), form)
	if ok {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateUser(ctx *gin.Context) {
	var (
		user model.User
		ok   bool
		msg  string
	)
	if middleware.IsAdmin(ctx) {
		var form f.UpdateUserForm
		if ok, msg = form.Bind(ctx); !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UpdateUserEventType)
		user = middleware.GetUser(ctx)
		ok, msg = service.UpdateUser(db.DB, user, form)
	} else {
		var form f.UpdateSelfForm
		if ok, msg = form.Bind(ctx); !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UpdateUserEventType)
		user = middleware.GetSelf(ctx).(model.User)
		ok, msg = service.UpdateSelf(db.DB, user, form)
	}
	if ok {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteUser(ctx *gin.Context) {
	var (
		tx  = db.DB.Begin()
		ok  bool
		msg string
	)
	if !middleware.IsAdmin(ctx) {
		var form f.DeleteSelfForm
		if ok, msg = form.Bind(ctx); !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.DeleteUserEventType)
		ok, msg = service.DeleteSelf(tx, middleware.GetSelf(ctx).(model.User), form)
	} else {
		ctx.Set(middleware.CTXEventTypeKey, model.DeleteUserEventType)
		ok, msg = db.InitUserRepo(tx).Delete(middleware.GetUser(ctx).ID)
	}
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
