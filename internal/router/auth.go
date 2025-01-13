package router

import (
	"RayWar/internal/db"
	"RayWar/internal/log"
	"RayWar/internal/middleware"
	"RayWar/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(ctx *gin.Context) {
	var registerForm RegisterForm
	trace := middleware.GetTraceID(ctx)
	if err := ctx.ShouldBind(&registerForm); err == nil {
		username, password, email := registerForm.Name, registerForm.Password, registerForm.Email
		user, ok, msg := db.CreateUser(username, password, email)
		if !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		if Token, err := utils.Generate(user.ID, user.Name, user.Type); err == nil {
			log.Logger.Infof("| %s | %s:%d register", trace, user.Name, user.ID)
			ctx.Writer.Header().Set("Authorization", "Bearer "+Token)
			ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": utils.TidyRetData(user, "password")[0]})
			return
		} else {
			msg = "UnknownError"
			log.Logger.Errorf("| %s | Generate token error: %s", trace, err)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
	} else {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}

func Login(ctx *gin.Context) {
	var loginForm LoginForm
	trace := middleware.GetTraceID(ctx)
	if err := ctx.ShouldBind(&loginForm); err == nil {
		username, password := loginForm.Name, loginForm.Password
		user, ok, msg := db.VerifyUser(username, password)
		if !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusUnauthorized, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		if Token, err := utils.Generate(user.ID, user.Name, user.Type); err == nil {
			log.Logger.Infof("| %s | %s:%d login", trace, user.Name, user.ID)
			ctx.Writer.Header().Set("Authorization", "Bearer "+Token)
			ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": utils.TidyRetData(user, "password")[0]})
			return
		} else {
			msg = "UnknownError"
			log.Logger.Errorf("| %s | Generate token error: %s", trace, err)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
	} else {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}
