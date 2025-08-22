package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WSAuth(ctx *gin.Context) {
	claims, err := utils.ParseToken(ctx.Query("token"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Unauthorized, "data": nil})
		return
	}
	DB := db.DB.WithContext(ctx)
	if claims.IsAdmin {
		admin, ok, msg := db.InitAdminRepo(DB).GetByID(claims.UserID)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		ctx.Set("IsAdmin", true)
		ctx.Set("Self", admin)
		ctx.Next()
	} else {
		user, ok, msg := db.InitUserRepo(DB).GetByID(claims.UserID)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		magic := ctx.Query("m")
		if !utils.CompareMagic(magic, claims.X) {
			go db.InitCheatRepo(db.DB.WithContext(ctx.Copy())).Create(db.CreateCheatOptions{
				UserID:  sql.Null[uint]{V: user.ID, Valid: true},
				Magic:   magic,
				IP:      ctx.ClientIP(),
				Reason:  fmt.Sprintf(model.DifferentTokenMagic, magic, claims.X),
				Type:    model.Suspicious,
				Checked: false,
			})
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Unauthorized, "data": nil})
			return
		}
		go db.InitDeviceRepo(db.DB.WithContext(ctx.Copy())).RecordDevice(user.ID, magic, ctx.ClientIP())
		if user.Banned {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Forbidden, "data": nil})
			return
		}
		ctx.Set("IsAdmin", false)
		ctx.Set("Self", user)
		ctx.Next()
	}
}
