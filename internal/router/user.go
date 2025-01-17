package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUser(ctx *gin.Context) {
	self, ok := ctx.Get("Self")
	if ok || self.(map[string]interface{})["Type"].(string) == "user" {
		user, ok, msg := db.GetUserByID(ctx, self.(map[string]interface{})["ID"].(uint), true)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": user})
		}
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
	}
}

func ChangePassword(ctx *gin.Context) {
	self, ok := ctx.Get("Self")
	if ok || self.(map[string]interface{})["Type"].(string) == "user" {
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
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden"})
	}
}

func UpdateUser(ctx *gin.Context) {
	self, ok := ctx.Get("Self")
	if ok || self.(map[string]interface{})["Type"].(string) == "user" {
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
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
	}
}
