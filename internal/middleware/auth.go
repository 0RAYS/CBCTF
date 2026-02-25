package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CheckAuth 是否登录, 用户是否被 ban, 记录设备
func CheckAuth(ctx *gin.Context) {
	auth := strings.Fields(ctx.GetHeader("Authorization"))
	if len(auth) != 2 || auth[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Response.Unauthorized})
		return
	}
	claims, err := utils.ParseToken(auth[1])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Response.Unauthorized})
		return
	}
	user, ret := db.InitUserRepo(db.DB).GetByID(claims.UserID)
	if !ret.OK {
		ctx.AbortWithStatusJSON(http.StatusOK, ret)
		return
	}
	magic := GetMagic(ctx)
	if !utils.CompareMagic(magic, claims.X) {
		if utils.CompareMagic(model.OauthLoginDeviceMagic, claims.X) {
			token, err := utils.GenerateToken(user.ID, user.Name, magic)
			if err != nil {
				log.Logger.Warningf("Failed to generate token: %s", err)
				ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
				return
			}
			ctx.Writer.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
		} else {
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
				prometheus.RecordCheatDetection(string(model.ReasonTypeTokenMagicType))
			}(contestIDL)
			ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Response.Unauthorized})
			return
		}
	}
	go db.InitDeviceRepo(db.DB).RecordDevice(db.CreateDeviceOptions{UserID: user.ID, Magic: magic})
	if user.Banned {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Response.Forbidden})
		return
	}
	ctx.Set("Self", user)
	ctx.Next()
}

// GetSelf 获取当前登录 admin 或 user
func GetSelf(ctx *gin.Context) model.User {
	self, ok := ctx.Get("Self")
	if !ok || self == nil {
		return model.User{}
	}
	return self.(model.User)
}

// CheckCaptain 检查是否为队伍队长, 要求 uri 中必须包含 contestID, admin 路由不能使用
func CheckCaptain(ctx *gin.Context) {
	team := GetTeam(ctx)
	if team.CaptainID != GetSelf(ctx).ID {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Response.Forbidden})
		return
	}
	ctx.Next()
}

// CheckVerified 检查邮箱是否已验证
func CheckVerified(ctx *gin.Context) {
	if !GetSelf(ctx).Verified {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Model.User.UnverifiedEmail})
		return
	}
	ctx.Next()
}

// CheckBanned 检查队伍是否被封禁
func CheckBanned(ctx *gin.Context) {
	team := GetTeam(ctx)
	if team.Banned {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Response.Forbidden})
		return
	}
	ctx.Next()
}
