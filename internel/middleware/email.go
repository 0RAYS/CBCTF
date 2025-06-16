package middleware

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckUnVerified(ctx *gin.Context) {
	user := GetSelf(ctx).(model.User)
	if user.Verified {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Next()
}
