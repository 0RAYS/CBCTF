package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUser(ctx *gin.Context) {
	self, _ := ctx.Get("Self")
	user, ok, msg := db.GetUserByID(ctx, self.(map[string]interface{})["ID"].(uint), true)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": user})
	}
}

func GetUsers(ctx *gin.Context) {
	var form GetUsersForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	users, count, ok, msg := db.GetUsers(ctx, form.Limit, form.Offset, false)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"count": count, "users": users}})
}

func ChangePassword(ctx *gin.Context) {
	self, _ := ctx.Get("Self")
	var form ChangePasswordForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest"})
		return
	}
	ok, msg := db.ChangePasswordUser(ctx, self.(map[string]interface{})["ID"].(uint), form.OldPassword, form.NewPassword)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
	}
}

func UpdateUser(ctx *gin.Context) {
	self, _ := ctx.Get("Self")
	var form UpdateForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	user, ok, msg := db.GetUserByID(ctx, self.(map[string]interface{})["ID"].(uint), false)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := utils.Form2Map(form)
	if user.Email != data["email"].(string) {
		if !db.IsUniqueEmail(data["email"].(string)) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "EmailExists", "data": nil})
			return
		}
		db.UpdateUser(ctx, self.(map[string]interface{})["ID"].(uint), map[string]interface{}{"verified": false})
	}
	if user.Name != data["name"].(string) && !db.IsUniqueName(data["name"].(string), model.User{}) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "UserNameExists", "data": nil})
		return
	}
	ok, msg = db.UpdateUser(ctx, self.(map[string]interface{})["ID"].(uint), data)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
	}
}
