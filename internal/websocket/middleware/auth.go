package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func WSAuth(ctx *gin.Context) {
	claims, err := utils.ParseToken(ctx.Query("token"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Unauthorized, "data": nil})
		return
	}
	if claims.IsAdmin {
		admin, ok, msg := db.InitAdminRepo(db.DB).GetByID(claims.UserID)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		ctx.Set("IsAdmin", true)
		ctx.Set("Self", admin)
		ctx.Next()
	} else {
		user, ok, msg := db.InitUserRepo(db.DB).GetByID(claims.UserID)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		magic := ctx.Query("m")
		if !utils.CompareMagic(magic, claims.X) {
			go db.InitCheatRepo(db.DB).Create(db.CreateCheatOptions{
				UserID:  sql.Null[uint]{V: user.ID, Valid: true},
				Magic:   magic,
				IP:      ctx.ClientIP(),
				Reason:  fmt.Sprintf(model.DifferentTokenMagic, magic, claims.X),
				Type:    model.Suspicious,
				Checked: false,
				Time:    time.Now(),
			})
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Unauthorized, "data": nil})
			return
		}
		go db.InitDeviceRepo(db.DB).RecordDevice(db.CreateDeviceOptions{UserID: user.ID, Magic: magic})
		if user.Banned {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Forbidden, "data": nil})
			return
		}
		ctx.Set("IsAdmin", false)
		ctx.Set("Self", user)
		ctx.Next()
	}
}
