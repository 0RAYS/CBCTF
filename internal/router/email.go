package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Verify(ctx *gin.Context) {
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
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}

func SendEmail(user model.User) (bool, string) {
	id := utils.UUID()
	token, err := utils.Generate(user.ID, user.Name, "email")
	if err != nil {
		return false, "UnknownError"
	}
	ok, msg := redis.SetEmailVerifyToken(user.ID, id)
	if !ok {
		return false, msg
	}
	if err := utils.SendVerifyEmail(user.Email, token, id); err != nil {
		return false, "SendEmailFailed"
	}
	return true, "Success"
}

func Activate(ctx *gin.Context) {
	user := middleware.GetSelf(ctx).(model.User)
	_, msg := SendEmail(user)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
