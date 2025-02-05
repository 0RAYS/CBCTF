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
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.ChangePasswordUser(tx, middleware.GetSelf(ctx).(model.User), form.OldPassword, form.NewPassword)
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
	tx := db.DB.WithContext(ctx).Begin()
	if email, ok := data["email"]; ok && email.(string) != user.Email {
		if !db.IsUniqueEmail(email.(string)) {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": "EmailExists", "data": nil})
			return
		}
		ok, msg = db.UpdateUser(tx, user.ID, map[string]interface{}{"verified": false})
		if !ok {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
	}
	if name, ok := data["name"]; ok && name.(string) != user.Name && !db.IsUniqueName(name.(string), model.User{}) {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": "UserNameExists", "data": nil})
		return
	}
	ok, msg := db.UpdateUser(tx, user.ID, data)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
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
	tx := db.DB.WithContext(ctx).Begin()
	// DeleteUser 需要嵌套预加载数据, 不可传入中间件保存的 ctx 数据
	ok, msg := db.DeleteUser(tx, ctx, userID)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func CreateUser(ctx *gin.Context) {
	var form constants.CreateUserForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	user, ok, msg := db.CreateUser(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
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
