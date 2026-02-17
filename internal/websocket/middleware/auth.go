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
			contestIDL, ret := db.InitContestRepo(db.DB).GetIDByUserID(user.ID)
			if !ret.OK {
				ctx.JSON(http.StatusOK, ret)
				return
			}
			go func(contestIDL []uint) {
				for _, contestID := range contestIDL {
					db.InitCheatRepo(db.DB).Create(db.CreateCheatOptions{
						ContestID:  contestID,
						Model:      model.CheatRefModel{user.ModelName(): {user.ID}},
						Magic:      magic,
						IP:         ctx.ClientIP(),
						Reason:     fmt.Sprintf(string(model.DifferentTokenMagicTmpl), magic, claims.X),
						ReasonType: model.ReasonTypeTokenMagicType,
						Type:       model.SuspiciousType,
						Checked:    false,
						Time:       time.Now(),
					})
				}
			}(contestIDL)
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
