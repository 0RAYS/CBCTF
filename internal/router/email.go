package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func VerifyEmail(ctx *gin.Context) {
	var form f.VerifyEmail
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	claims, err := utils.Parse(form.Token)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "BadEmailVerifyToken", "data": nil})
		return
	}
	id, ok := redis.GetEmailVerifyToken(claims.UserID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": id, "data": nil})
		return
	}
	if form.ID == id {
		redis.DelEmailVerifyToken(claims.UserID)
		tx := db.DB.WithContext(ctx).Begin()
		ok, msg := db.UpdateUser(tx, claims.UserID, map[string]interface{}{"verified": true})
		if !ok {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		tx.Commit()
	}
	ctx.Redirect(http.StatusMovedPermanently, config.Env.Frontend)
}

func ActivateEmail(ctx *gin.Context) {
	user := middleware.GetSelf(ctx).(model.User)
	token, err := utils.Generate(user.ID, user.Name, "email")
	id := utils.RandomString()
	if err != nil || !redis.SetEmailVerifyToken(user.ID, id) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "UnknownError", "data": nil})
		return
	}
	if err := utils.SendVerifyEmail(user.Email, token, id); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "SendEmailError", "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}
