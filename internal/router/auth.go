package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(ctx *gin.Context) {
	var form RegisterForm
	trace := middleware.GetTraceID(ctx)
	if err := ctx.ShouldBind(&form); err == nil {
		username, password, email := form.Name, form.Password, form.Email
		user, ok, msg := db.CreateUser(ctx, username, password, email)
		if !ok {
			log.Logger.Infof("%s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		if Token, err := utils.Generate(user.ID, user.Name, "user"); err == nil {
			log.Logger.Infof("%s | %s:%d register", trace, user.Name, user.ID)
			ctx.Writer.Header().Set("Authorization", "Bearer "+Token)
			ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": user})
			return
		} else {
			msg = "UnknownError"
			log.Logger.Errorf("%s | Generate token error: %s", trace, err)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
	} else {
		log.Logger.Infof("%s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}

func Login(ctx *gin.Context) {
	var form LoginForm
	trace := middleware.GetTraceID(ctx)
	if err := ctx.ShouldBind(&form); err == nil {
		username, password := form.Name, form.Password
		user, ok, msg := db.VerifyUser(ctx, username, password)
		if !ok {
			log.Logger.Infof("%s | %s", trace, msg)
			ctx.JSONP(http.StatusUnauthorized, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		if Token, err := utils.Generate(user.ID, user.Name, "user"); err == nil {
			log.Logger.Infof("%s | %s:%d login", trace, user.Name, user.ID)
			ctx.Writer.Header().Set("Authorization", "Bearer "+Token)
			ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": user})
			return
		} else {
			msg = "UnknownError"
			log.Logger.Errorf("%s | Generate token error: %s", trace, err)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
	} else {
		log.Logger.Infof("%s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}

func AdminLogin(ctx *gin.Context) {
	var form LoginForm
	trace := middleware.GetTraceID(ctx)
	if err := ctx.ShouldBind(&form); err == nil {
		username, password := form.Name, form.Password
		admin, ok, msg := db.VerifyAdmin(ctx, username, password)
		if !ok {
			log.Logger.Infof("%s | %s", trace, msg)
			ctx.JSONP(http.StatusUnauthorized, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		if Token, err := utils.Generate(admin.ID, admin.Name, "admin"); err == nil {
			log.Logger.Infof("%s | %s:%d login", trace, admin.Name, admin.ID)
			ctx.Writer.Header().Set("Authorization", "Bearer "+Token)
			ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": admin})
			return
		} else {
			msg = "UnknownError"
			log.Logger.Errorf("%s | Generate token error: %s", trace, err)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
	} else {
		log.Logger.Infof("%s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}
