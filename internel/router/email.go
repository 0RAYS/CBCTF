package router

import (
	"CBCTF/internel/config"
	f "CBCTF/internel/form"
	"CBCTF/internel/log"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	"CBCTF/internel/redis"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"CBCTF/internel/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func VerifyEmail(ctx *gin.Context) {
	var form f.VerifyEmail
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.VerifyEmail(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}

func ActivateEmail(ctx *gin.Context) {
	user := middleware.GetSelf(ctx).(model.User)
	_, msg := SendEmail(user)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func SendEmail(user model.User) (bool, string) {
	if strings.ToLower(config.Env.Gin.Mode) == "debug" {
		return true, "DebugMode"
	}
	id := utils.UUID()
	token, err := utils.Generate(user.ID, user.Name, "email")
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		return false, "UnknownError"
	}
	ok, msg := redis.SetEmailVerifyToken(user.ID, id)
	if !ok {
		return false, msg
	}
	if err := utils.SendVerifyEmail(user.Email, token, id); err != nil {
		return false, "SendEmailError"
	}
	return true, "Success"
}
