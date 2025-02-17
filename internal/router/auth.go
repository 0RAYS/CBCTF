package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(ctx *gin.Context) {
	var form f.RegisterForm
	trace := middleware.GetTraceID(ctx)
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	user, ok, msg := db.CreateUser(tx, f.CreateUserForm{Name: form.Name, Password: form.Password, Email: form.Email})
	if !ok {
		tx.Rollback()
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	if Token, err := utils.Generate(user.ID, user.Name, "user"); err == nil {
		log.Logger.Infof("%s | %s:%d register", trace, user.Name, user.ID)
		ctx.Writer.Header().Set("Authorization", "Bearer "+Token)
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": &user})
		return
	} else {
		msg = "UnknownError"
		log.Logger.Warningf("%s | Generate token error: %s", trace, err)
		ctx.JSONP(http.StatusInternalServerError, gin.H{"msg": msg, "data": nil})
		return
	}
}

func Login(ctx *gin.Context) {
	var form f.LoginForm
	trace := middleware.GetTraceID(ctx)
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	username, password := form.Name, form.Password
	user, ok, msg := db.VerifyUser(db.DB.WithContext(ctx), username, password)
	if !ok {
		ctx.JSONP(http.StatusUnauthorized, gin.H{"msg": msg, "data": nil})
		return
	}
	if Token, err := utils.Generate(user.ID, user.Name, "user"); err == nil {
		log.Logger.Infof("%s | %s:%d login", trace, user.Name, user.ID)
		ctx.Writer.Header().Set("Authorization", "Bearer "+Token)
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": &user})
		return
	} else {
		msg = "UnknownError"
		log.Logger.Warningf("%s | Generate token error: %s", trace, err)
		ctx.JSONP(http.StatusInternalServerError, gin.H{"msg": msg, "data": nil})
		return
	}
}

func AdminLogin(ctx *gin.Context) {
	var form f.LoginForm
	trace := middleware.GetTraceID(ctx)
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	username, password := form.Name, form.Password
	admin, ok, msg := db.VerifyAdmin(db.DB.WithContext(ctx), username, password)
	if !ok {
		ctx.JSONP(http.StatusUnauthorized, gin.H{"msg": msg, "data": nil})
		return
	}
	if Token, err := utils.Generate(admin.ID, admin.Name, "admin"); err == nil {
		log.Logger.Infof("%s | %s:%d login", trace, admin.Name, admin.ID)
		ctx.Writer.Header().Set("Authorization", "Bearer "+Token)
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": &admin})
		return
	} else {
		msg = "UnknownError"
		log.Logger.Warningf("%s | Generate token error: %s", trace, err)
		ctx.JSONP(http.StatusInternalServerError, gin.H{"msg": msg, "data": nil})
		return
	}
}
