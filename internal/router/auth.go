package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/oauth"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var form f.RegisterForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.RegisterEventType)
	tx := db.DB.WithContext(ctx).Begin()
	user, ok, msg := service.CreateUser(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if ok, msg = service.SendEmail(user); !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	token, err := utils.GenerateToken(user.ID, user.Name, false, middleware.GetMagic(ctx))
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	ctx.Set("IsAdmin", false)
	ctx.Set("Self", user)
	log.Logger.Infof("%s:%d register", user.Name, user.ID)
	ctx.Writer.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	prometheus.UpdateUserRegisterMetrics(oauth.LocalProvider)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.RegisterResp(user, false)})
}

func Login(ctx *gin.Context) {
	var form f.LoginForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.LoginEventType)
	user, ok, msg := service.VerifyUser(db.DB.WithContext(ctx), form)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	token, err := utils.GenerateToken(user.ID, user.Name, false, middleware.GetMagic(ctx))
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	ctx.Set("IsAdmin", false)
	ctx.Set("Self", user)
	log.Logger.Infof("%s:%d login", user.Name, user.ID)
	ctx.Writer.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	prometheus.UpdateUserLoginMetrics(oauth.LocalProvider)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.LoginResp(user, false)})
}

func AdminLogin(ctx *gin.Context) {
	var form f.LoginForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.LoginEventType)
	admin, ok, msg := service.VerifyAdmin(db.DB.WithContext(ctx), form)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	token, err := utils.GenerateToken(admin.ID, admin.Name, true, "admin")
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	ctx.Set("IsAdmin", true)
	ctx.Set("Self", admin)
	log.Logger.Infof("%s:%d login", admin.Name, admin.ID)
	ctx.Writer.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.GetAdminResp(admin)})
}
