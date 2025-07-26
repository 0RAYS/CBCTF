package router

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	db "CBCTF/internal/repo"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
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
	users, count, ok, msg := db.InitUserRepo(db.DB.WithContext(ctx)).List(form.Limit, form.Offset, db.GetOptions{
		Preloads: map[string]db.GetOptions{
			"Teams":    {},
			"Contests": {},
		},
	})
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
	tx := db.DB.WithContext(ctx).Begin()
	user, ok, msg := service.AdminCreateUser(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.GetUserResp(user, true)})
}

func ChangePwd(ctx *gin.Context) {
	var form f.ChangePasswordForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.ChangeUserPwd(tx, middleware.GetSelf(ctx).(model.User), form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateUser(ctx *gin.Context) {
	var (
		user model.User
		ok   bool
		msg  string
		tx   *gorm.DB
	)
	if middleware.IsAdmin(ctx) {
		var form f.UpdateUserForm
		if ok, msg := form.Bind(ctx); !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		user = middleware.GetUser(ctx)
		tx = db.DB.WithContext(ctx).Begin()
		ok, msg = service.UpdateUser(tx, user, form)
	} else {
		var form f.UpdateSelfForm
		if ok, msg = form.Bind(ctx); !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		user = middleware.GetSelf(ctx).(model.User)
		tx = db.DB.WithContext(ctx).Begin()
		ok, msg = service.UpdateSelf(tx, user, form)
	}
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteUser(ctx *gin.Context) {
	var (
		tx  = db.DB.WithContext(ctx).Begin()
		ok  bool
		msg string
	)
	if !middleware.IsAdmin(ctx) {
		var form f.DeleteSelfForm
		if ok, msg = form.Bind(ctx); !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		ok, msg = service.DeleteSelf(tx, middleware.GetSelf(ctx).(model.User), form)
	} else {
		ok, msg = db.InitUserRepo(tx).Delete(middleware.GetUser(ctx).ID)
	}
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
