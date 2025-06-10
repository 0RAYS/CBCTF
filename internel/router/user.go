package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetUser(ctx *gin.Context) {
	all := middleware.GetRole(ctx) == "admin"
	var user model.User
	if !all {
		user = middleware.GetSelf(ctx).(model.User)
	} else {
		user = middleware.GetUser(ctx)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": resp.GetUserResp(user, all)})
}

func GetUsers(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 5
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	users, count, ok, msg := db.InitUserRepo(db.DB.WithContext(ctx)).List(form.Limit, form.Offset, "Teams", "Contests")
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
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
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
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest})
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
	if middleware.GetRole(ctx) == "admin" {
		var form f.UpdateUserForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
			return
		}
		user = middleware.GetUser(ctx)
		tx = db.DB.WithContext(ctx).Begin()
		ok, msg = service.UpdateUser(tx, user, form)
	} else if middleware.GetRole(ctx) == "user" {
		var form f.UpdateSelfForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
			return
		}
		user = middleware.GetSelf(ctx).(model.User)
		tx = db.DB.WithContext(ctx).Begin()
		ok, msg = service.UpdateSelf(tx, user, form)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Forbidden, "data": nil})
		return
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
	if middleware.GetRole(ctx) != "admin" {
		var form f.DeleteSelfForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest})
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
