package router

import (
	f "CBCTF/internel/form"
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
	var (
		user model.User
		all  bool
	)

	if middleware.GetRole(ctx) != "admin" {
		user = middleware.GetSelf(ctx).(model.User)
		all = false
	} else {
		user = middleware.GetUser(ctx)
		all = true
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": resp.GetUserResp(user, all)})
}

func ChangePwd(ctx *gin.Context) {
	var form f.ChangePasswordForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest"})
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
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		user = middleware.GetUser(ctx)
		tx = db.DB.WithContext(ctx).Begin()
		ok, msg = service.UpdateUser(tx, user, form)
	} else if middleware.GetRole(ctx) == "user" {
		var form f.UpdateSelfForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		user = middleware.GetSelf(ctx).(model.User)
		tx = db.DB.WithContext(ctx).Begin()
		ok, msg = service.UpdateSelf(tx, user, form)
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
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
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest"})
			return
		}
		ok, msg = service.DeleteSelf(tx, middleware.GetSelf(ctx).(model.User), form)
	} else {
		ok, msg = service.DeleteUser(tx, middleware.GetUser(ctx))
	}
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
