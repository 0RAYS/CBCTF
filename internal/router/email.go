package router

import (
	"CBCTF/internal/config"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	db "CBCTF/internal/repo"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyEmail(ctx *gin.Context) {
	var form f.VerifyEmail
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.VerifyEmailEventType)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.VerifyEmail(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.Redirect(http.StatusTemporaryRedirect, config.Env.Backend)
}

func ActivateEmail(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.ActivateEmailEventType)
	user := middleware.GetSelf(ctx).(model.User)
	if user.Verified {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.AlreadyVerifiedEmail, "data": nil})
		return
	}
	_, msg := service.SendEmail(user)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
