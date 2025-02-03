package router

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUser(ctx *gin.Context) {
	if middleware.GetRole(ctx) != "admin" {
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": middleware.GetSelf(ctx).(model.User)})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": middleware.GetUser(ctx)})
}

func ChangePassword(ctx *gin.Context) {
	var form constants.ChangePasswordForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest"})
		return
	}
	_, msg := db.ChangePasswordUser(ctx, middleware.GetSelf(ctx).(model.User), form.OldPassword, form.NewPassword)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})

}

func UpdateUser(ctx *gin.Context) {
	var (
		user model.User
		msg  string
		data map[string]interface{}
	)
	if middleware.GetRole(ctx) == "admin" {
		var form constants.UpdateUserForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		user = middleware.GetUser(ctx)
		data = utils.Form2Map(form)
		if password, ok := data["password"]; ok && password.(string) != "" {
			data["password"] = utils.HashPassword(password.(string))
		} else {
			data["password"] = user.Password
		}
	} else if middleware.GetRole(ctx) == "user" {
		var form constants.UpdateSelfForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		user = middleware.GetSelf(ctx).(model.User)
		data = utils.Form2Map(form)
		data["password"] = user.Password
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
		return
	}
	if email, ok := data["email"]; ok && email.(string) != user.Email {
		if !db.IsUniqueEmail(email.(string)) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "EmailExists", "data": nil})
			return
		}
		db.UpdateUser(ctx, user.ID, map[string]interface{}{"verified": false})
	}
	if name, ok := data["name"]; ok && name.(string) != user.Name && !db.IsUniqueName(name.(string), model.User{}) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "UserNameExists", "data": nil})
		return
	}
	_, msg = db.UpdateUser(ctx, user.ID, data)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteUser(ctx *gin.Context) {
	var userID uint
	if middleware.GetRole(ctx) != "admin" {
		var form constants.DeleteSelfForm
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest"})
			return
		}
		if !utils.CompareHashAndPassword(middleware.GetSelf(ctx).(model.User).Password, form.Password) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "PasswordError", "data": nil})
			return
		}
		userID = middleware.GetSelfID(ctx)
	} else {
		userID = middleware.GetUser(ctx).ID
	}
	// DeleteUser 需要嵌套预加载数据, 不可传入中间件保存的 ctx 数据
	_, msg := db.DeleteUser(ctx, userID)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func CreateUser(ctx *gin.Context) {
	var form constants.CreateUserForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	user, ok, msg := db.CreateUser(ctx, form)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": user})
}

func GetUsers(ctx *gin.Context) {
	var form constants.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	users, count, ok, msg := db.GetUsers(ctx, form.Limit, form.Offset, true)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "users": users}})
}
