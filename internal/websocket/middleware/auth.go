package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func WSAuth(ctx *gin.Context) {
	claims, err := utils.ParseToken(ctx.Query("token"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Request.Unauthorized})
		return
	}
	if claims.IsAdmin {
		admin, ret := db.InitAdminRepo(db.DB).GetByID(claims.UserID)
		if !ret.OK {
			ctx.AbortWithStatusJSON(http.StatusOK, ret)
			return
		}
		ctx.Set("IsAdmin", true)
		ctx.Set("Self", admin)
		ctx.Next()
	} else {
		user, ret := db.InitUserRepo(db.DB).GetByID(claims.UserID)
		if !ret.OK {
			ctx.AbortWithStatusJSON(http.StatusOK, ret)
			return
		}
		magic := GetMagic(ctx)
		if !utils.CompareMagic(magic, claims.X) {
			go db.InitCheatRepo(db.DB).Create(db.CreateCheatOptions{
				Model:   map[string]uint{user.ModelName(): user.ID},
				Magic:   magic,
				IP:      ctx.ClientIP(),
				Reason:  fmt.Sprintf(model.DifferentTokenMagic, magic, claims.X),
				Type:    model.Suspicious,
				Checked: false,
				Time:    time.Now(),
			})
			ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Request.Unauthorized})
			return
		}
		go db.InitDeviceRepo(db.DB).RecordDevice(db.CreateDeviceOptions{UserID: user.ID, Magic: magic})
		if user.Banned {
			ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Request.Forbidden})
			return
		}
		ctx.Set("IsAdmin", false)
		ctx.Set("Self", user)
		ctx.Next()
	}
}
