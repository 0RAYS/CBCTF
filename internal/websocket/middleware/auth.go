package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// parseWSProtocols parses the Sec-WebSocket-Protocol header value into a map.
// The browser sends protocols as a comma-separated list, e.g. "Bearer, <token>, Magic, <magic>".
// We treat consecutive pairs as key→value.
func parseWSProtocols(header string) map[string]string {
	result := make(map[string]string)
	parts := strings.Split(header, ",")
	for i := 0; i+1 < len(parts); i += 2 {
		key := strings.TrimSpace(parts[i])
		val := strings.TrimSpace(parts[i+1])
		result[key] = val
	}
	return result
}

func WSAuth(ctx *gin.Context) {
	protocols := parseWSProtocols(ctx.Request.Header.Get("Sec-Websocket-Protocol"))
	tokenValue := protocols["Bearer"]
	claims, err := utils.ParseToken(tokenValue)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Request.Unauthorized})
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
	ctx.Set("Self", user)
	ctx.Next()

}
