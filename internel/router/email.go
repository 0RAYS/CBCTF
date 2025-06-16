package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	"CBCTF/internel/redis"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"CBCTF/internel/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func VerifyEmail(ctx *gin.Context) {
	var form f.VerifyEmail
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
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
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}

func ActivateEmail(ctx *gin.Context) {
	user := middleware.GetSelf(ctx).(model.User)
	_, msg := SendEmail(user)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func SendEmail(user model.User) (bool, string) {
	id := utils.UUID()
	token, err := utils.Generate(user.ID, user.Name, "email", "email")
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		return false, i18n.UnknownError
	}
	ok, msg := redis.SetEmailVerifyToken(user.ID, id)
	if !ok {
		return false, msg
	}
	if err = utils.SendVerifyEmail(user.Email, token, id); err != nil {
		log.Logger.Warningf("Failed to send mail: %s", err)
		return false, i18n.SendEmailError
	}
	return true, i18n.Success
}
