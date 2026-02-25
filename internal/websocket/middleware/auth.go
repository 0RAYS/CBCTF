package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// parseWSToken extracts the JWT token and magic from the Sec-WebSocket-Protocol header.
// Frontend encodes each value as base64url (no padding) and prefixes with "token-" / "magic-"
// to stay within the valid HTTP token character set required by RFC 6455.
func parseWSToken(header string) (token, magic string) {
	for _, part := range strings.Split(header, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "token-") {
			decoded, err := base64.RawURLEncoding.DecodeString(part[len("token-"):])
			if err == nil {
				token = string(decoded)
			}
		} else if strings.HasPrefix(part, "magic-") {
			decoded, err := base64.RawURLEncoding.DecodeString(part[len("magic-"):])
			if err == nil {
				magic = string(decoded)
			}
		}
	}
	return
}

func WSAuth(ctx *gin.Context) {
	tokenValue, _ := parseWSToken(ctx.Request.Header.Get("Sec-Websocket-Protocol"))
	claims, err := utils.ParseToken(tokenValue)
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
		contestIDL, ret := db.InitContestRepo(db.DB).GetIDByUserID(user.ID)
		if !ret.OK {
			ctx.JSON(http.StatusOK, ret)
			return
		}
		ip := ctx.ClientIP()
		go func(contestIDL []uint) {
			for _, contestID := range contestIDL {
				db.InitCheatRepo(db.DB).Create(db.CreateCheatOptions{
					ContestID:  contestID,
					Model:      model.CheatRefModel{user.ModelName(): {user.ID}},
					Magic:      magic,
					IP:         ip,
					Reason:     fmt.Sprintf(string(model.DifferentTokenMagicTmpl), magic, claims.X),
					ReasonType: model.ReasonTypeTokenMagicType,
					Type:       model.SuspiciousType,
					Checked:    false,
					Time:       time.Now(),
				})
			}
		}(contestIDL)
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Response.Unauthorized})
		return
	}
	go db.InitDeviceRepo(db.DB).RecordDevice(db.CreateDeviceOptions{UserID: user.ID, Magic: magic})
	if user.Banned {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Response.Forbidden})
		return
	}
	ctx.Set("Self", user)
	ctx.Next()
}
