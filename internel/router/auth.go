package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/log"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"CBCTF/internel/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(ctx *gin.Context) {
	var form f.RegisterForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	user, ok, msg := service.CreateUser(tx, f.RegisterForm{
		Name:     form.Name,
		Password: form.Password,
		Email:    form.Email,
	})
	if !ok {
		tx.Rollback()
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if ok, msg = SendEmail(user); !ok {
		tmp := tx.Rollback()
		log.Logger.Warning(tmp)
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	token, err := utils.Generate(user.ID, user.Name, "user")
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		ctx.JSONP(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
		return
	}
	log.Logger.Infof("%s:%d register", user.Name, user.ID)
	ctx.Writer.Header().Set("Authorization", "Bearer "+token)
	ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": resp.RegisterResp(user)})
	return
}

func Login(ctx *gin.Context) {
	var form f.LoginForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	user, ok, msg := service.VerifyUser(db.DB.WithContext(ctx), form)
	if !ok {
		ctx.JSONP(http.StatusUnauthorized, gin.H{"msg": msg, "data": nil})
		return
	}
	if Token, err := utils.Generate(user.ID, user.Name, "user"); err == nil {
		log.Logger.Infof("%s:%d login", user.Name, user.ID)
		ctx.Writer.Header().Set("Authorization", "Bearer "+Token)
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": resp.LoginResp(user)})
		return
	} else {
		msg = "UnknownError"
		log.Logger.Warningf("Generate token error: %s", err)
		ctx.JSONP(http.StatusInternalServerError, gin.H{"msg": msg, "data": nil})
		return
	}
}

func AdminLogin(ctx *gin.Context) {
	var form f.LoginForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	admin, ok, msg := service.VerifyAdmin(db.DB.WithContext(ctx), form)
	if !ok {
		ctx.JSONP(http.StatusUnauthorized, gin.H{"msg": msg, "data": nil})
		return
	}
	if Token, err := utils.Generate(admin.ID, admin.Name, "admin"); err == nil {
		log.Logger.Infof("%s:%d login", admin.Name, admin.ID)
		ctx.Writer.Header().Set("Authorization", "Bearer "+Token)
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": &admin})
		return
	} else {
		msg = "UnknownError"
		log.Logger.Warningf("Generate token error: %s", err)
		ctx.JSONP(http.StatusInternalServerError, gin.H{"msg": msg, "data": nil})
		return
	}
}
