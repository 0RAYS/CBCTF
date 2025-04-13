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
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	user, ok, msg := service.CreateUser(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if ok, msg = SendEmail(user); !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	token, err := utils.Generate(user.ID, user.Name, "user")
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
		return
	}
	log.Logger.Infof("%s:%d register", user.Name, user.ID)
	ctx.Writer.Header().Set("Authorization", "Bearer "+token)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.RegisterResp(user, false)})
	return
}

func Login(ctx *gin.Context) {
	var form f.LoginForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	user, ok, msg := service.VerifyUser(db.DB.WithContext(ctx), form)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": msg, "data": nil})
		return
	}
	if Token, err := utils.Generate(user.ID, user.Name, "user"); err == nil {
		log.Logger.Infof("%s:%d login", user.Name, user.ID)
		ctx.Writer.Header().Set("Authorization", "Bearer "+Token)
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.LoginResp(user, false)})
		return
	} else {
		msg = "UnknownError"
		log.Logger.Warningf("Failed to generate token: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": msg, "data": nil})
		return
	}
}

func AdminLogin(ctx *gin.Context) {
	var form f.LoginForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	admin, ok, msg := service.VerifyAdmin(db.DB.WithContext(ctx), form)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": msg, "data": nil})
		return
	}
	if Token, err := utils.Generate(admin.ID, admin.Name, "admin"); err == nil {
		log.Logger.Infof("%s:%d login", admin.Name, admin.ID)
		ctx.Writer.Header().Set("Authorization", "Bearer "+Token)
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.GetAdminResp(admin)})
		return
	} else {
		msg = "UnknownError"
		log.Logger.Warningf("Failed to generate token: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": msg, "data": nil})
		return
	}
}
