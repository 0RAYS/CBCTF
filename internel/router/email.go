package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
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
	_, msg := service.SendEmail(user)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
